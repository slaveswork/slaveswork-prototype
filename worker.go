package main

import "github.com/Equanox/gotron"

type Worker struct {
	window  *gotron.BrowserWindow
	address Address
}

func newWorker(w *gotron.BrowserWindow) *Worker {
	return &Worker{
		window: w,
	}
}
