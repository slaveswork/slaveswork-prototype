package main

import (
	"log"
	"net"
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



	return &addr, listener // return address object and listener for http handler.
}

func getIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("func : getIP\nError : ", err)
	}

	var ip string

	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	return ip
}