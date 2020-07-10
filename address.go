package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
)

type Address struct {
	IP   string `json:"ip"`   // IP address for application
	Port string `json:"port"` // Port number for application --> can change to integer
}

func newAddress() (Address, net.Listener) {
	addr := Address{}                        // application address(IP + Port)
	listener, err := net.Listen("tcp", ":0") // for finding unused Port number.
	if err != nil {
		log.Fatal("func : newAddress\n", err)
	}

	addr.IP = getIP()                                             // get IP address
	addr.Port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port) // get Port number

	return addr, listener // return address object and listener for http handler.
}

// get IP address using 'net' package.
func getIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("func : getIP\n", err)
	}

	var ip string // return value

	// valid IPv4 address(ignore Loop back address...127.0.0.1)
	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	return ip
}

func (a *Address) generateToken() string {
	s := a.IP + a.Port // IP -> string, Port -> string | ex) "127.0.0.180" = "127.0.0.1"(IP) + "80"(Port)

	// make SHA1 hash value --> can change another hash function
	h := sha1.New()
	h.Write([]byte(s)) // s -> "{ip}{port}" string value
	bs := h.Sum(nil)

	token := fmt.Sprintf("%x", bs) // byte slice to string

	// SHA1 hash value is too long to be used by users.
	return token[:12] // return only 12 characters(string).
}

func (a *Address) generateHostAddress(path string) (hostAddress string) {
	u := url.URL{
		Scheme: "http",
		Host:   a.IP + ":" + a.Port,
		Path:   path,
	}

	hostAddress = u.String()
	return
}