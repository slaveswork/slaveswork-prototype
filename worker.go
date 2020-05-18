package main

import (
	"encoding/json"
	"github.com/Equanox/gotron"
	"log"
	"net/http"
)

type Worker struct {
	window  *gotron.BrowserWindow
	address Address
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
	w.address, _ = newAddress()
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
	receivedMessage.Body = body

	// "app.connect.device" Unmarshal JSON message
	if err := json.Unmarshal(bin, &receivedMessage); err != nil {
		log.Fatal("func : sendConnectionRequest\nError : ", err)
	}
	checkJSON(receivedMessage)

	// 'IP', 'Port', 'Token' --> is in message's Body.
	url := body.Address.generateHostAddress(body.Token)
}
