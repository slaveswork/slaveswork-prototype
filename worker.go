package main

import (
	"bytes"
	"encoding/json"
	"github.com/Equanox/gotron"
	"log"
	"net/http"
)

type Worker struct {
	window      *gotron.BrowserWindow `json:"-"`
	hostAddress Address               `json:"-"`

	Id      int     `json:"id"`
	Address Address `json:"address"`
}

func newWorker(w *gotron.BrowserWindow) *Worker {
	return &Worker{
		window: w,
	}
}

func (w *Worker) run() {
	w.init()
	w.gotronMessageHandler()
}

func (w *Worker) init() {
	w.Address, _ = newAddress()
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

	respId := struct {
		Id int `json:"id"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&respId); err != nil {
		log.Fatal("func : sendConnectionRequest\n", err)
	}

	w.Id = respId.Id
}
