package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Equanox/gotron"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Worker struct {
	window      *gotron.BrowserWindow `json:"-"` // to ignore in marshaling
	hostAddress Address               `json:"-"` // to ignore in marshaling
	ci          chan int              `json:"-"`

	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Method  string  `json:"method"`
	Address Address `json:"address"`

	// Worker status -> ignore
	Status      string      `json:"-"`
	T           chan Tile   `json:"-"`
	FilePath    chan string `json:"-"`
	Config		Config
	BlenderPath string `json:"-"`
}

func newWorker(w *gotron.BrowserWindow) *Worker {
	return &Worker{
		window: w, // for worker's window.
		ci:     make(chan int),
		T:           make(chan Tile),
		FilePath:    make(chan string),
	}
}

func (w *Worker) run() {
	w.init()
	w.gotronMessageHandler()
	w.sendWorkerStatus()
	w.sendBlenderPath()
	w.sendHostIp()
	go w.httpMessageHandler()
	go w.renderTileWithBlender()
}

func (w *Worker) init() {
	w.Name, _ = os.Hostname()
	w.Address, _ = newAddress() // initialize worker's address.
	w.Config = UnmarshalConfig()
}

func (w *Worker) gotronMessageHandler() {
	w.window.On(&gotron.Event{Event: "app.connect.device"}, w.sendConnectionRequest) // connection between worker and host.
	w.window.On(&gotron.Event{Event: "app.blender.path"}, w.receiveBlenderPath) // blender.exe path
}

func (w *Worker) httpMessageHandler() {
	http.HandleFunc("/render/resource", w.receiveRenderResource)

	http.ListenAndServe(":" + w.Address.Port, nil)
}

func (w *Worker) sendConnectionRequest(bin []byte) {
	var receivedMessage GotronMessage // for Unmarshal window message
	body := struct {
		Address
		Token string `json:token` // Token...
	}{}
	receivedMessage.Body = &body

	// "app.connect.device" Unmarshal JSON message
	if err := json.Unmarshal(bin, &receivedMessage); err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}
	w.hostAddress = body.Address // Initialize host address.
	w.Method = "add" // first send "add" worker to host.

	w.Config.HostIp = body.Address.IP
	w.Config.SaveConfig()

	// Marshaling body for connection request.
	requestBody, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}

	// send request to host.
	url := body.Address.generateHostAddress(body.Token) // make host url --> format : http://{ip}:{port}/{token}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}

	// resp's body will have this worker's Id(for Host's worker management).
	defer resp.Body.Close()

	respId := struct { // temporary struct for receiving response body.
		Id int `json:"id"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&respId); err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}

	w.Id = respId.Id // set up worker's Id for host's management.
	w.ci <- respId.Id
}

func (w *Worker) sendWorkerStatus() {
	<-w.ci      // wait function until setting worker's Id.
	close(w.ci) // close channel... we don't need that anymore.

	url := w.hostAddress.generateHostAddress("status") // Host address.
	w.Method = "update" // from now on, send "update" method to host.

	for range time.Tick(time.Second * 5) { // repeat sending this worker's status every 5 seconds.
		requestBody, err := json.MarshalIndent(w, "", "    ")
		if err != nil {
			log.Fatal("func : sendWorkerStatus\n", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatal("func : sendWorkerStatus\n", err)
		}

		respBytes, _ := ioutil.ReadAll(resp.Body) // read response from host.
		fmt.Println(string(respBytes)) // printing for validation.
		resp.Body.Close() // close response body.
	}
}

func (w *Worker) sendBlenderPath() {
	message := GotronMessage{
		Event: &gotron.Event{Event: "window.blender.path"},
	}
	message.Body = struct { // temporary struct for sending token message.
		BlenderPath string `json:"blenderPath"`
	}{
		BlenderPath: w.Config.BlenderPath, // initialize value.
	}

	checkJSON(message) // Printing Message for validation.
	w.window.Send(message)
}

func (w *Worker) sendHostIp() {
	message := GotronMessage{
		Event: &gotron.Event{Event: "window.send.hostIp"},
	}
	message.Body = struct { // temporary struct for sending token message.
		HostIp string `json:"hostIp"`
	}{
		HostIp: w.Config.HostIp, // initialize value.
	}

	checkJSON(message) // Printing Message for validation.
	w.window.Send(message)
}

func (w *Worker) receiveBlenderPath(bin []byte) {
	var message GotronMessage
	body := struct {
		blenderPath string `json:"blenderPath"`
	}{}
	message.Body = &body

	if err := json.Unmarshal(bin, &message); err != nil {
		log.Fatal("func : receiveBlenderPath\n", err)
	}

	w.Config.BlenderPath = body.blenderPath
	w.Config.SaveConfig()
}

func (w *Worker) receiveRenderResource(rw http.ResponseWriter, r *http.Request) {
	t := Tile {
		Index: r.Header.Get("index"),
		Xmin: r.Header.Get("xmin"),
		Ymin: r.Header.Get("ymin"),
		Xmax: r.Header.Get("xmax"),
		Ymax: r.Header.Get("ymax"),
		Frame: r.Header.Get("fram"),
	}

	w.T <- t

	file, header, err := r.FormFile("blend-file")
	if err != nil {
		log.Fatal("func : receiveRenderResource\n", err)
	}

	defer file.Close()

	path := os.TempDir() + "/" + header.Filename
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(rw, "Unable to create the file for writing.")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(rw, err)
	}

	w.FilePath <- path
}

func (w *Worker) renderTileWithBlender() {
	var blendFile string
	var tile Tile

	for {
		blendFile = <- w.FilePath // *.blend file path
		tile = <- w.T

		cmd := exec.Command(w.Config.BlenderPath,
			"-b", blendFile,
			"-F", "EXR",
			"-Y",
			"-noaudio",
			"-E", "CYCLES",
			"-P", "worker.py",
			"--", tile.Xmin, tile.Ymin, tile.Xmax, tile.Ymax, tile.Frame)

		stdout, err := cmd.StdoutPipe()
		stderr, err := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			log.Fatal("func : renderTileWithBlender\n", err)
		}

		go copyOutput(stdout)
		go copyOutput(stderr)
		cmd.Wait()
	}
}

func copyOutput(r io.Reader) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
}