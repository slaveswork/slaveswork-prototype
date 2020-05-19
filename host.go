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
	h.httpMessageHandler(listener)
}

func (h *Host) init() (listener net.Listener) {
	// 'h.address' is already created by 'newHost' function
	h.address, listener = newAddress() // initialize Host's IP address and Port number.
	// send network status to 'Host' window.
	h.send("window.network.status")
	h.token = h.address.generateToken() // generate Token for Worker's connection.

	go h.handleWorkers() // register & unregister Worker.

	return // listener --> net.Listener for http function handler.
}

// this function manage workers for host.
func (h *Host) handleWorkers() {
	i := 0 // indexing workers

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
	h.window.On(&gotron.Event{Event: "app.generate.token"}, func(bin []byte) { h.send("window.send.token") })
}

func (h *Host) httpMessageHandler(listenser net.Listener) {
	http.HandleFunc("/"+h.token, h.sendConnectionResponse)

	http.Serve(listenser, nil)
}

func (h *Host) sendConnectionResponse(w http.ResponseWriter, r *http.Request) {
	var worker Worker // for unmarshal request
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		log.Fatal("func : sendConnectionResponse\n", err)
	}

	h.register <- &worker // register worker for management
	respId := struct { // make struct for response
		Id int `json:"id"`
	}{
		Id: <-h.index,
	}

	// make response body for ResponseWriter
	respBody, err := json.MarshalIndent(respId, "", "    ")
	if err != nil {
		log.Fatal("func : sendConnectionResponse\n", err)
	}

	// set response Header and send to worker
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}
