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
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Token string `json:"token"`
}

type WindowNetworkStatusMessage struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type WindowSendTokenMessage struct {
	Token string `json:"token"`
}

func (c *AppConnectionDeviceMessage) MakeHostAddress() string {
	return fmt.Sprintf("%s:%s", c.IP, c.Port)
}