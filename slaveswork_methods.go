package main

import (
	"crypto/sha1"
	"fmt"
	"net"
	"strings"
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