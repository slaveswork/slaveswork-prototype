package message

import "github.com/Equanox/gotron"

type AppWebSocketEvent struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}

type AppConnectionDeviceEvent struct {
	*gotron.Event                       `json:"event"`
	Message *AppConnectionDeviceMessage `json:"message"`
}

type AppConnectionDeviceMessage struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Token string `json:"token"`
}

type WindowNetworkStatusEvent struct {
	*gotron.Event                       `json:"event"`
	Message *WindowNetworkStatusMessage `json:"message"`
}

type WindowNetworkStatusMessage struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type WindowSendTokenEvent struct {
	*gotron.Event                   `json:"event"`
	Message *WindowSendTokenMessage `json:"message"`
}

type WindowSendTokenMessage struct {
	Token string `json:"token"`
}

type HostConnectWorkerMessage struct {
	Token string `json:"token"`
}