package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Usage: %s (IP address)\n", os.Args[0])
		os.Exit(1)
	}

	// 获取输入的ip地址
	address := os.Args[1]

	// 解析IP
	ip := net.ParseIP(address)
	if ip == nil {
		fmt.Fprintf(os.Stderr, "invalid address")
		os.Exit(1)
	}

	// 获取 IP 地址默认的子网掩码
	defaultMask := ip.DefaultMask()
	fmt.Printf("Subnet mask is: %s\n", defaultMask)

	// 获取主机所在的网络地址
	network := ip.Mask(defaultMask)
	fmt.Printf("Network address is: %s\n", network)

	// 获取掩码位数和掩码总长度
	ones, bits := defaultMask.Size()
	fmt.Printf("Mask bits: %d, Total bits: %d\n", ones, bits)
}
