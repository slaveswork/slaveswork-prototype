package main

import (
	"encoding/json"
	"fmt"
	"github.com/Equanox/gotron"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type Worker struct {
	window *gotron.BrowserWindow
	conn   *websocket.Conn
}

/*
TODO make a package for these Event Message structures.
slaveswork_back --- main.go
                 +- worker.go
                 +- host.go
                 |
                 +- event --- event.go
*/

type connectionEvent struct {
	*gotron.Event
	Message *connectionMessage `json:"message"`
}

type connectionMessage struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Token string `json:"token"`
}

type hostConnectionEvent struct {
	Event string `json:"event"`
	Message *connectionMessage `json:"message"`
}

func newWorker(w *gotron.BrowserWindow) *Worker {
	return &Worker {
		window: w,
		conn:   nil,
	}
}

func (w *Worker) run() {
	var connEvent connectionEvent

	w.window.On(&gotron.Event{Event: "app.connect.device"}, func(bin []byte) {
		// make connection request for host.
		fmt.Println(fmt.Sprintf("%x", string(bin)))
		if err := json.Unmarshal(bin, &connEvent); err != nil {
			log.Fatal("App connect JSON parsing : ", err)
		}

		// connection message's port is missing. But host do 'not' need that.
		// TODO IP can change to MAC address.
		// TODO Change JSON to another.
		connRequestEvent := &hostConnectionEvent{
			Event:   "host.connect.worker",
			Message: &connectionMessage{
				IP:    getIPAddress(),           // worker's IP address.
				Token: connEvent.Message.Token,  // host's Token for validation.
			},
		}

		u := url.URL{
			Scheme: "ws",
			Host:   connEvent.makeHostAddress(),
			Path:   "/connectWorker",
		}
		fmt.Printf("Connecting to %s\n", u.String())

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial : ", err)
		}

		defer conn.Close()

		if err := conn.WriteJSON(connRequestEvent); err != nil {
			log.Fatal("Connection Device WriteJSON error : ", err)
		}
		w.conn = conn
	})
}

func (c *connectionEvent) makeHostAddress() string {
	return fmt.Sprintf("%s:%s", c.Message.IP, c.Message.Port)
}