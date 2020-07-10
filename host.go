package main

import (
	"encoding/json"
	"fmt"
	"github.com/Equanox/gotron"
	"log"
	"net"
	"net/http"
)

type Host struct {
	window *gotron.BrowserWindow // for communication with window.

	address    Address // for host IP address and Port number.
	token      string
	workers    map[int]*Worker
	index      chan int
	register   chan *Worker
	unregister chan int
}

func newHost(w *gotron.BrowserWindow) *Host {
	return &Host{
		window: w,

		workers:    make(map[int]*Worker),
		index:      make(chan int),
		register:   make(chan *Worker),
		unregister: make(chan int),
	}
}

func (h *Host) run() {
	listener := h.init()
	h.gotronMessageHandler()
	go h.httpMessageHandler(listener)
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

	// handlers for blender
	http.HandleFunc("/running/check", h.receiveRunningCheck)
	http.HandleFunc("/task/resource", ReceiveTaskResource)

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

	h.workers[worker.Id] = &worker

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
		tile.prettyPrint()
	}
}