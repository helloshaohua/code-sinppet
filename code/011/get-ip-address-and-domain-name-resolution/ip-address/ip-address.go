package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s (IP address)\n", os.Args[0])
		os.Exit(1)
	}

	// Get the input ip address
	address := os.Args[1]

	// Parse IP
	ip := net.ParseIP(address)
	if ip != nil {
		fmt.Printf("The address is %s.\n", ip.String())
	} else {
		fmt.Println("invalid address")
	}
	os.Exit(0)
}
