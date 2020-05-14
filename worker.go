package main

import (
	"encoding/json"
	"fmt"
	"github.com/Equanox/gotron"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"github.com/slaveswork/slaveswork-prototype/message"
)

type Worker struct {
	window *gotron.BrowserWindow
	conn   *websocket.Conn
}

func newWorker(w *gotron.BrowserWindow) *Worker {
	return &Worker{
		window: w,
		conn:   nil,
	}
}

func (w *Worker) run() {
	var appConnectionDeviceEvent message.AppConnectionDeviceMessage

	w.window.On(&gotron.Event{Event: "app.connect.device"}, func(bin []byte) {
		// make connection request for host.
		fmt.Println(fmt.Sprintf("%x", string(bin)))
		if err := json.Unmarshal(bin, &appConnectionDeviceEvent); err != nil {
			log.Fatal("JSON parsing : ", err)
		}

		u := url.URL{
			Scheme: "ws",
			Host:   appConnectionDeviceEvent.Body.MakeHostAddress(),
			Path:   fmt.Sprintf("/%s", appConnectionDeviceEvent.Body.Token),
		}
		fmt.Printf("Connecting to %s\n", u.String())

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("Connection Websocket Dial : ", err)
		}

		w.conn = conn
	})
}
