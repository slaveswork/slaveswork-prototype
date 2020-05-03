package message

import (
	"fmt"
	"github.com/Equanox/gotron"
)

type GotronMessage struct {
	*gotron.Event    `json:"event"`
	Body interface{} `json:"body"`
}

type AppMessage struct {
	Event string      `json:"event"`
 	Body  interface{} `json:"body"`
}

type AppConnectionDeviceMessage struct {
	*gotron.Event
	Body *AppConnectionDeviceBody `json:"body"`
}

type AppConnectionDeviceBody struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Token string `json:"token"`
}

type WindowNetworkStatusMessage struct {
	*gotron.Event                 `json:"event"`
	Body *WindowNetworkStatusBody `json:"body"`
}

type WindowNetworkStatusBody struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type WindowSendTokenMessage struct {
	*gotron.Event             `json:"event"`
	Body *WindowSendTokenBody `json:"body"`
}

type WindowSendTokenBody struct {
	Token string `json:"token"`
}

func (c *AppConnectionDeviceBody) MakeHostAddress() string {
	return fmt.Sprintf("%s:%s", c.IP, c.Port)
}