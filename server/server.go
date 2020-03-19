package server

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/asticode/go-astilectron"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/multiformats/go-multiaddr"
)

// ClientInfo is JSON format for client cpu usage.
type clientInfo struct {
	Name  string    `json:"name"`
	Usage []float64 `json:"usage"`
}

var window *astilectron.Window

// LaunchServer run host application
func LaunchServer(w *astilectron.Window, c chan string) {
	window = w
	sourcePort := 4444

	var r io.Reader = rand.Reader

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", sourcePort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	if err != nil {
		panic(err)
	}

	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	// Only applies on the receiving side.
	host.SetStreamHandler("/chat/1.0.0", handleStream)

	// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
	var port string
	for _, la := range host.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		panic("was not able to find actual local port")
	}

	fmt.Printf("Run './chat -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
	fmt.Println("You can replace 127.0.0.1 with public IP as well.")
	fmt.Printf("\nWaiting for incoming connection\n\n")

	destination := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%s", port, host.ID().Pretty())
	c <- destination

	// Hang forever
	<-make(chan struct{})
}

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readCPU(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func readCPU(rw *bufio.ReadWriter) {
	// while true loop
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}

		if str != "\n" {
			res := clientInfo{}
			json.Unmarshal([]byte(str), &res)
			// Send message to electron application
			window.SendMessage(str, func(m *astilectron.EventMessage) {
				// Unmarshal
				var s string
				m.Unmarshal(&s)

				// Process message
				fmt.Println(fmt.Sprintf("Host Application Updates %s for %s", s, res.Name))
			})
		}
	}
}
