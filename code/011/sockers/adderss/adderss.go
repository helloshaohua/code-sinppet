package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s IP-address", os.Args[0])
		os.Exit(1)
	}

	address := os.Args[1]

	addr := net.ParseIP(address)

	if addr == nil {
		fmt.Println("Invalid address")
	} else {
		fmt.Println(addr.String())
	}
	os.Exit(0)

	net.Dial()
}
