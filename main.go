package main

import (
	"github.com/Equanox/gotron"
)

func main() {
	// Create a new browser window instance
	window := app()

	// Start the browser window.
	// This will establish a golang <=> nodejs bridge using websockets,
	// to control ElectronBrowserWindow with our window object.
	done, err := window.Start()
	if err != nil {
		panic(err)
	}

	// Open dev tools must be used after window.Start
	window.OpenDevTools()

	processEvent(window)

	// Wait for the application to close
	<- done
}

func app() *gotron.BrowserWindow {
	window, err := gotron.New("ui/dist")
	if err != nil {
		panic(err)
	}

	// Alter default window size and window title.
	window.WindowOptions.Width = 1432
	window.WindowOptions.Height = 700
	window.WindowOptions.Title = "Slave's work"

	return window
}

func processEvent(w *gotron.BrowserWindow) {
	w.On(&gotron.Event{Event: "app.host.start"}, func(bin []byte) {
		host := newHost(w)
		go host.run()
	})

	w.On(&gotron.Event{Event: "app.worker.start"}, func(bin []byte) {
		worker := newWorker(w)
		go worker.run()
	})
}