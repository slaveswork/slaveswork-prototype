package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Equanox/gotron"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
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
	Port  uint16 `json:"port"`
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
		if err := json.Unmarshal(bin, &connEvent); err != nil {
			panic(err)
		}

		// connection message's port is missing. But host do 'not' need that.
		// TODO IP can change to MAC address.
		connRequestEvent := &hostConnectionEvent{
			Event:   "",
			Message: &connectionMessage{
				IP:    getIPAddress(),           // worker's IP address.
				Token: connEvent.Message.Token,  // host's Token for validation.
			},
		}

		u := url.URL{
			Scheme: "http",
			Host:   connEvent.makeHostAddress(),
			Path:   "/host.connect.worker",
		}

		// Send http request to host with POST method.
		eventMessage, _ := json.Marshal(connRequestEvent)
		eventMessageBuff := bytes.NewBuffer(eventMessage)

		resp, err := http.Post(u.String(), "application/json", eventMessageBuff)
		if err != nil {
			panic(err)
		}

		result, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		log.Println(result)
	})
}

func (c *connectionEvent) makeHostAddress() string {
	return fmt.Sprintf("%s:%d", c.Message.IP, c.Message.Port)
}