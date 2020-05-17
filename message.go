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
	Event *gotron.Event `json:"event"` // is event name.
	Body  interface{}   `json:"body"`  // is actual content.
}

// Host's message sender : send message to window.
func (h *Host) send(e string, b interface{}) {
	message := GotronMessage{
		Event: &gotron.Event{Event: e},
	}

	switch e {
	case "window.network.status":
		message.Body, _ = message.Body.(Address) // convert interface{} to 'Address'(in main/address.go)
		message.Body = b.(Address) // put parameter 'b' on message's Body. --> have to convert 'interface{}' to 'Address'(in main/address.go)
	}

	prettyJson, err := json.MarshalIndent(message, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate JSON", err)
	}
	fmt.Printf("%s\n", string(prettyJson))

	h.window.Send(message)
}
