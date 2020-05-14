package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/Equanox/gotron"
	"github.com/gorilla/websocket"
	"github.com/slaveswork/slaveswork-prototype/message"
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
	ip         string
	port       string
	token      string
}

func newHost(w *gotron.BrowserWindow) *Host {
	return &Host{
		window:     w,
		workers:    make(map[*Worker]bool),
		register:   make(chan *Worker),
		unregister: make(chan *Worker),
		ip:         "127.0.0.1",
		port:       "80",
	}
}

func (h *Host) run() {
	listener := h.init()

	h.handleWindowMessages()

	h.handleHttpMessages(listener)
}

func (h *Host) init() net.Listener {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	h.ip = getIPAddress()
	h.port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	h.token = makeToken(h.ip)

	h.window.Send(message.WindowNetworkStatusMessage{
		Event: &gotron.Event{Event: "window.network.status"},
		Body: &message.WindowNetworkStatusBody{
			IP:   h.ip,
			Port: h.port,
		},
	})

	go func() {
		for {
			select {
			case worker := <-h.register:
				h.workers[worker] = true
				fmt.Println("Registered Worker...")
			case worker := <-h.unregister:
				if _, ok := h.workers[worker]; ok {
					delete(h.workers, worker)
				}
			}
		}
	}()

	return listener
}

func (h *Host) handleWindowMessages() {
	h.window.On(&gotron.Event{Event: "app.generate.token"}, func(bin []byte) {
		h.window.Send(message.WindowSendTokenMessage{
			Event: &gotron.Event{Event: "window.send.token"},
			Body:  &message.WindowSendTokenBody{
				Token: h.token,
			},
		})
	})
}

func (h *Host) handleHttpMessages(listener net.Listener) {
	// make handler for worker's connection request.
	// Register handler.
	// Later, It should be replace with 'http.HandleFunc(pattern, connectionRequest)'
	// connectionRequest <-- function name
	pattern := fmt.Sprintf("/%s", h.token)
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
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