package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/Equanox/gotron"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
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
	Event   *gotron.Event         `json:"event"`
	Message *networkStatusMessage `json:"message"`
}

func (n networkStatusEvent) EventString() string {
	panic("implement me")
}

type networkStatusMessage struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type sendTokenEvent struct {
	Event   *gotron.Event     `json:"event"`
	Message *sendTokenMessage `json:"message"`
}

func (s sendTokenEvent) EventString() string {
	panic("implement me")
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
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	ip    := getIPAddress()
	port  := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	token := makeToken(ip)

	messgae := networkStatusEvent {
		Event:   &gotron.Event{Event: "window.network.status"},
		Message: &networkStatusMessage{
			IP: ip,
			Port: port,
		},
	}

	h.window.Send(messgae)

	tokenMessage := sendTokenEvent {
		Event:   &gotron.Event{Event: "window.send.token"},
		Message: &sendTokenMessage{
			Token: token,
		},
	}

	h.window.On(&gotron.Event{Event: "app.generate.token"}, func(bin []byte) {
		h.window.Send(tokenMessage)
	})

	// waiting for channel...
	go func() {
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
	}()

	// make handler for worker's connection request.
	// Register handler.
	http.HandleFunc("/connectWorker", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received Request for connection with worker")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		// create worker for management
		worker := &Worker{window: nil, conn: conn}
		h.register <- worker // register worker
	})

	// Start HandleFunc
	http.Serve(listener, nil)
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
		if addrs, err := interf.Addrs(); err == nil {
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
	token := fmt.Sprintf("%x", bs)

	return token[:12]
}