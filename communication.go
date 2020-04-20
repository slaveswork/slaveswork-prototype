package main

import (
	"github.com/Equanox/gotron"
	"log"
)

type CustomEvent struct {
	Event   *gotron.Event `json:"event"`
	Message string        `json:"message"`
}

func (c CustomEvent) EventString() string {
	panic("implement me")
}

func ReceiveMessage(window *gotron.BrowserWindow, channel string, callback func(bin []byte)) {
	window.On(&gotron.Event{Event: channel}, callback)
}

func SendMessage(window *gotron.BrowserWindow, channel string, msg string) {
	event := CustomEvent{
		Event:   &gotron.Event{Event: channel},
		Message: msg,
	}

	if err := window.Send(event); err != nil {
		log.Println("send Error !", err)
	}
}
