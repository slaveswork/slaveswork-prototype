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
	var appConnectionDeviceEvent message.AppConnectionDeviceEvent

	w.window.On(&gotron.Event{Event: "app.connect.device"}, func(bin []byte) {
		// make connection request for host.
		fmt.Println(fmt.Sprintf("%x", string(bin)))
		if err := json.Unmarshal(bin, &appConnectionDeviceEvent); err != nil {
			log.Fatal("JSON parsing : ", err)
		}

		// connection message's port is missing. But host do 'not' need that.
		// TODO IP can change to MAC address.
		// TODO Change JSON to another.
		connRequestMessage := &message.HostConnectWorkerMessage{
			Token: appConnectionDeviceEvent.Message.Token, // Token for validation.
		}

		bMessage, err := json.Marshal(connRequestMessage)
		if err != nil {
			log.Fatal("Connection Request Marshal : ", err)
		}

		connRequestEvent := &message.AppWebSocketEvent{
			Event:   "host.connect.worker",
			Message: string(bMessage),
		}

		u := url.URL{
			Scheme: "ws",
			Host:   appConnectionDeviceEvent.Message.MakeHostAddress(),
			Path:   "/connectWorker",
		}
		fmt.Printf("Connecting to %s\n", u.String())

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("Connection Websocket Dial : ", err)
		}

		defer conn.Close()

		if err := conn.WriteJSON(connRequestEvent); err != nil {
			log.Fatal("Connection Device WriteJSON error : ", err)
		}
		w.conn = conn
	})
}
