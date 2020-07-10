package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Equanox/gotron"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
}

func newWorker(w *gotron.BrowserWindow) *Worker {
	return &Worker{
		window: w, // for worker's window.
		ci:     make(chan int),
	}
}

func (w *Worker) run() {
	w.init()
	w.gotronMessageHandler()
	w.sendWorkerStatus()
}

func (w *Worker) init() {
	w.Name, _ = os.Hostname()
	w.Address, _ = newAddress() // initialize worker's address.
}

func (w *Worker) gotronMessageHandler() {
	w.window.On(&gotron.Event{Event: "app.connect.device"}, w.sendConnectionRequest) // connection between worker and host.
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
