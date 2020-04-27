package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/Equanox/gotron"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"github.com/slaveswork/slaveswork-prototype/message"
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

	networkStatusMessage := message.GotronMessage{
		Event: &gotron.Event{Event: "window.network.status"},
		Body:  &message.WindowNetworkStatusMessage{
			IP:   ip,
			Port: port,
		},
	}

	h.window.Send(networkStatusMessage)

	sendTokenMessage := message.GotronMessage{
		Event: &gotron.Event{Event: "window.send.token"},
		Body:  &message.WindowSendTokenMessage{
			Token: token,
		},
	}

	h.window.On(&gotron.Event{Event: "app.generate.token"}, func(bin []byte) {
		h.window.Send(sendTokenMessage)
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
					worker.conn.Close()
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