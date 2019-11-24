package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s network address\n", os.Args[0])
		os.Exit(1)
	}

	// 网络类型
	network := os.Args[1]

	// 域名或IP
	address := os.Args[2]

	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.ResolveTCPAddr error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "network: %s, address: %s\n", tcpAddr.Network(), tcpAddr.String())
	os.Exit(0)
}
