package main

import "github.com/Equanox/gotron"

type Host struct {
	window  *gotron.BrowserWindow // for communication with window.
	address Address               // for host IP address and Port number.
}

func newHost(w *gotron.BrowserWindow) *Host {
	return &Host{
		window: w,
	}
}

func (h *Host) run() {

}

func (h *Host) init() {

}