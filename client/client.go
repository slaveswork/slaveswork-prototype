package client

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/shirou/gopsutil/cpu"
)

// ClientInfo is JSON format for client cpu usage.
type clientInfo struct {
	Name  string    `json:"name"`
	Usage []float64 `json:"usage"`
}

// LaunchClient run client application
func LaunchClient(c chan string, dest string) {
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

	fmt.Println("This node's multiaddresses:")
	for _, la := range host.Addrs() {
		fmt.Printf(" - %v\n", la)
	}
	fmt.Println()

	// Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		log.Fatalln(err)
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalln(err)
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
	if err != nil {
		panic(err)
	}

	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	go writeCPU(rw)

	c <- "Client Start"

	// Hang forever.
	select {}
}

// Send message about CPU usage to server
func writeCPU(rw *bufio.ReadWriter) {
	for {
		cpuUsages, err := cpu.Percent(time.Duration(1)*time.Second, true)
		hostName, _ := os.Hostname()

		data := &clientInfo{
			Name:  hostName,
			Usage: cpuUsages}

		strData, err := json.Marshal(data)

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", string(strData)))
		rw.Flush()
	}
}
