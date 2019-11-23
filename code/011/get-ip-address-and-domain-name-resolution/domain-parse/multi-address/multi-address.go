package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s hostname\n", os.Args[0])
		os.Exit(1)
	}

	// 获取命令输入域名
	hostname := os.Args[1]

	// 通过域名获取 IP
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LookupHost error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Lookup host address is: %s\n", strings.Join(addrs, ", "))
	os.Exit(0)
}
