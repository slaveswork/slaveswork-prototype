package main

import (
	"crypto/sha1"
	"github.com/Equanox/gotron"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{}

type Host struct {
	window     *gotron.BrowserWindow
	workers    map[*Worker]bool
	register   chan *Worker
	unregister chan *Worker
}

/*
TODO make a package for these Event Message structures.
slaveswork_back --- main.go
                 +- worker.go
                 +- host.go
                 |
                 +- event --- event.go
*/

type networkStatusEvent struct {
	*gotron.Event
	Message *networkStatusMessage `json:"message"`
}

type networkStatusMessage struct {
	IP string `json:"ip"`
}

type sendTokenEvent struct {
	*gotron.Event
	Message *sendTokenMessage `json:"message"`
}

type sendTokenMessage struct {
	Token string `json:"token"`
}

func newHost(w *gotron.BrowserWindow) *Host {
	return &Host {
		window:     w,
		workers:    make(map[*Worker]bool),
		register:   make(chan *Worker),
		unregister: make(chan *Worker),
	}
}

func (h *Host) run() {
	ip := getIPAddress()
	token := makeToken(ip)

	h.window.Send(&networkStatusEvent{
		Event:   &gotron.Event{Event: "window.network.status"},
		Message: &networkStatusMessage{
			IP: ip,
		},
	})

	h.window.On(&gotron.Event{Event: "app.generate.token"}, func(bin []byte) {
		h.window.Send(&sendTokenEvent{
			Event:   &gotron.Event{Event: "window.send.token"},
			Message: &sendTokenMessage{
				Token: token,
			},
		})
	})

	// make handler for worker's connection request.
	// Register handler.
	http.HandleFunc("/host.connect.worker", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		worker := &Worker{window: nil, conn: conn}
		h.register <- worker
	})

	// waiting for channel...
	for {
		select {
		case worker := <-h.register:
			h.workers[worker] = true
		case worker := <-h.unregister:
			if _, ok := h.workers[worker]; ok {
				delete(h.workers, worker)
			}
		}
	}
}

/*
Now, this function is 'Public' method.
But worker also need this function for finding worker's IP address.
So make another package for methods like this.
 */
func getIPAddress() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	var currentIP string

	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				currentIP = ipNet.IP.String()
			}
		}
	}

	return currentIP
}

func makeToken(ip string) string {
	var currentNetworkHardwareName string

	interfaces, _ := net.Interfaces()

	for _, interf := range interfaces {
		if addrs, err := interf.Addrs(); err != nil {
			for _, addr := range addrs {
				if strings.Contains(addr.String(), ip) {
					currentNetworkHardwareName = interf.Name
				}
			}
		}
	}

	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)
	if err != nil {
		panic(err)
	}

	macAddress := netInterface.HardwareAddr
	h := sha1.New()
	h.Write([]byte(macAddress))
	bs := h.Sum(nil)

	// TODO Should be modified
	return string(bs[:12])
}