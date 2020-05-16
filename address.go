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