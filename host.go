package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Equanox/gotron"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Host struct {
	window *gotron.BrowserWindow // for communication with window.

	address    Address // for host IP address and Port number.
	token      string
	workers    map[int]*Worker
	index      chan int
	register   chan *Worker
	unregister chan int

	filePath   chan string
	tiles      []Tile
	freeWorker chan bool
}

func newHost(w *gotron.BrowserWindow) *Host {
	return &Host{
		window: w,

		workers:    make(map[int]*Worker),
		index:      make(chan int),
		register:   make(chan *Worker),
		unregister: make(chan int),

		filePath:   make(chan string),
		freeWorker: make(chan bool, 1),
	}
}

func (h *Host) run() {
	listener := h.init()
	h.gotronMessageHandler()
	go h.httpMessageHandler(listener)

	worker := newWorker(h.window)
	worker.run()

	var message GotronMessage
	body := struct {
		Address
		Token string `json:token`
	}{
		Address{h.address.IP, h.address.Port},
		h.token,
	}
	message.Body = &body
	request, _ := json.Marshal(message)
	worker.sendConnectionRequest(request)
}

func (h *Host) init() (listener net.Listener) {
	// 'h.address' is already created by 'newHost' function
	h.address, listener = newAddress() // initialize Host's IP address and Port number.
	// send network status to 'Host' window.
	h.send("window.network.status", nil)
	h.token = h.address.generateToken() // generate Token for Worker's connection.

	go h.handleWorkers() // register & unregister Worker.

	return // listener --> net.Listener for http function handler.
}

// this function manage workers for host.
func (h *Host) handleWorkers() {
	i := 1 // indexing workers

	for {
		select {
		case worker := <-h.register: // register worker
			fmt.Println("Register Worker...")
			h.workers[i] = worker // append worker at host's map
			h.index <- i // return current worker's index --> go to responseWriter
			i += 1 // change index for next worker

		case unregisteredIndex := <-h.unregister: // unregister worker
			delete(h.workers, unregisteredIndex) // remove current worker from host's map
		}
	}
}

func (h *Host) gotronMessageHandler() {
	h.window.On(&gotron.Event{Event: "app.generate.token"}, func(bin []byte) { h.send("window.send.token", nil) })
}

func (h *Host) httpMessageHandler(listener net.Listener) {
	http.HandleFunc("/"+h.token, h.receiveConnectionRequest)
	http.HandleFunc("/status", h.receiveWorkerStatus)
	http.HandleFunc("/task/tiles", h.receiveTilesInfo)
	http.HandleFunc("/task/result", h.receiveTaskResult)

	// handlers for blender
	http.HandleFunc("/running/check", h.receiveRunningCheck)
	http.HandleFunc("/task/resource", h.ReceiveTaskResource)

	http.Serve(listener, nil)
}

func (h *Host) receiveConnectionRequest(w http.ResponseWriter, r *http.Request) {
	var worker Worker // for unmarshal request
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		log.Fatal("func : receiveConnectionRequest\n", err)
	}

	h.register <- &worker // register worker for management

	respId := struct { // make struct for response
		Id int `json:"id"`
	}{
		Id: <-h.index,
	}

	(&worker).Id = respId.Id
	(&worker).Status = "WAITING"
	// Add worker status at host's window. ( Method : "Add" )
	h.send("window.device.status", &worker)

	// make response body for ResponseWriter
	respBody, err := json.MarshalIndent(respId, "", "    ")
	if err != nil {
		log.Fatal("func : receiveConnectionRequest\n", err)
	}

	// set response Header and send to worker
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}

func (h *Host) receiveWorkerStatus(w http.ResponseWriter, r *http.Request) {
	var worker Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		log.Fatal("func : receiveWorkerStatus\n", err)
	}

	h.workers[worker.Id].Cpu = worker.Cpu
	h.workers[worker.Id].Memory = worker.Memory

	// Update worker status at host's window. ( Method : "Update" )
	h.send("window.device.status", &worker)
}

func (h *Host) receiveRunningCheck(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Blender Add-On : checking host running...")
	w.WriteHeader(http.StatusOK) // return status code 200
	// return 404 when host is not running
}

func (h *Host) receiveTilesInfo(w http.ResponseWriter, r *http.Request) {
	var tiles []Tile
	if err :=  json.NewDecoder(r.Body).Decode(&tiles); err != nil{
		log.Fatal("{error}, tiles unmarshall error ")
	}

	for _, tile := range tiles {
		tile.prettyPrint() // printing tiles information
	}

	h.tiles = tiles
	h.DoRender() // tiles : tile information, path : blender file path
}

func (h *Host) ReceiveTaskResource(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("func : ReceiveTaskResource\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	fmt.Println("blen : ", string(reqBody[0:4]))
	if blen := string(reqBody[:4]); strings.Compare(blen, "blen") != 0 {
		log.Fatal("func : ReceiveTaskResource\n", "This is not a task resource request.")
	}
	buffer := bytes.NewBuffer(reqBody[4:8])
	fileSize, err := binary.ReadVarint(buffer)
	fmt.Println("size : ", fileSize)

	if err != nil {
		log.Fatal("func : ReceiveTaskResource\n", "This is not a number.")
	}

	//if fileSize != len(reqBody[8:]) {
	//	log.Fatal("func : ReceiveTaskResource\n", "An error occurred while receiving the file.")
	//}

	dir := filepath.Join(os.TempDir(), "main.blend")
	err = ioutil.WriteFile(dir, reqBody[8:], 0644)
	if err != nil {
		log.Fatal("func : ReceiveTaskResource\n", "File Writing Error.")
	}
	fmt.Println("dir:", dir)

	h.filePath <- dir
}

func (h *Host) receiveTaskResult(w http.ResponseWriter, r *http.Request) {
	wId, _ := strconv.Atoi(r.Header.Get("Worker"))
	tIdx := r.Header.Get("Tile")

	file, header, err := r.FormFile("blend-file")

	defer file.Close()

	path := filepath.Join(os.TempDir(), header.Filename)
	out, err := os.Create(path)
	if err != nil {
		log.Fatal("func : receiveTaskResult\n", err)
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(w, err)
	}

	h.workers[wId].Status = "WAITING"
	for i, _ := range h.tiles {
		if h.tiles[i].Index == tIdx {
			h.tiles[i].Success = true
		}
	}

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		log.Fatal("func : receiveTaskResult\n", err)
	}

	req.Header.Add("filename", path)
	req.Header.Add("index", tIdx)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal("func : receiveTaskResult\n", err)
	}

	for len(h.freeWorker) > 0 {
		<- h.freeWorker
	}
	h.freeWorker <- true
}