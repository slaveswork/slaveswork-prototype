package main

import (
	"net"
)

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