// UDP Client端设计
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port\n", os.Args[0])
		os.Exit(1)
	}

	// 服务IP和端口号
	service := os.Args[1]

	// 获取 UDPAddr 结构体指针
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	// 连接到UDP网络服务器
	udpConn, err := net.DialUDP("udp", nil, udpAddr)

	// 发送数据到UDP网络服务器
	_, err = udpConn.Write([]byte("hello server"))
	checkError(err)

	// 读取UDP网络服务器响应数据
	bytes := make([]byte, 512)
	n, addr, err := udpConn.ReadFromUDP(bytes)
	checkError(err)
	fmt.Printf("server address: %s, recevied response data: %s\n", addr.String(), string(bytes[0:n]))

	udpConn.Close()
	os.Exit(0)
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
