package main

import (
	"log"
	"net"
	"strconv"
)

type Address struct {
	ip   string // IP address for application
	port string // Port number for application --> can change to integer
}

func newAddress() (*Address, net.Listener) {
	addr := Address{} // application address(IP + Port)
	listener, err := net.Listen("tcp", ":0") // for finding unused Port number.
	if err != nil {
		log.Fatal("func : newAddress\nError : ", err)
	}

	addr.ip = getIP() // get IP address
	addr.port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port) // get Port number

	return &addr, listener // return address object and listener for http handler.
}

// get IP address using 'net' package.
func getIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("func : getIP\nError : ", err)
	}

	var ip string // return value

	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	return ip
}