package main

import (
	"github.com/Equanox/gotron"
	"net"
)

type Host struct {
	window  *gotron.BrowserWindow // for communication with window.
	address Address               // for host IP address and Port number.
	token   string
}

func newHost(w *gotron.BrowserWindow) *Host {
	return &Host{
		window: w,
	}
}

func (h *Host) run() {
	_ = h.init()
	h.gotronMessageHandler()
}

func (h *Host) init() (listener net.Listener) {
	// 'h.address' is already created by 'newHost' function
	h.address, listener = newAddress() // initialize Host's IP address and Port number.
	// send network status to 'Host' window.
	h.send("window.network.status")
	h.token = h.address.generateToken() // generate Token for Worker's connection.
	return // listener --> net.Listener for http function handler.
}

func (h *Host) gotronMessageHandler() {
	h.window.On(&gotron.Event{Event: "app.generate.token"}, func(bin []byte) {h.send("window.send.token")})
}