package main

import (
	"encoding/json"
	"fmt"
	"github.com/Equanox/gotron"
	"log"
)

func (g GotronMessage) EventString() string {
	panic("implement me")
}

type GotronMessage struct {
	*gotron.Event // is event name.
	Body  interface{} `json:"body"` // is actual content.
}

// Host's message sender : send message to window.
func (h *Host) send(e string) {
	message := GotronMessage{
		Event: &gotron.Event{Event: e},
	}

	switch e {
	case "window.network.status":
		message.Body, _ = message.Body.(Address) // convert interface{} to 'Address'(in main/address.go)
		message.Body = h.address                 // put parameter 'b' on message's Body. --> have to convert 'interface{}' to 'Address'(in main/address.go)
	case "window.send.token":
		message.Body = struct { // temporary struct for sending token message.
			Token string `json:"token"`
		}{
			Token: h.token, // initialize value.
		}
	}

	checkJSON(message) // Printing Message for validation.
	h.window.Send(message)
}

// Pretty printing JSON message.
func checkJSON(message GotronMessage) {
	prettyJson, err := json.MarshalIndent(message, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate JSON", err)
	}
	fmt.Printf("%s\n", string(prettyJson))
}
