// UDP Server 端设计
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 获取 UDPAddr 结构体指针
	udpAddr, err := net.ResolveUDPAddr("udp", ":8826")
	checkError(err)

	// 监听服务端口
	udpConn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	// 循环处理每一个客户端
	for {
		handleClient(udpConn)
	}
}

func handleClient(udpConn *net.UDPConn) {
	bytes := make([]byte, 512)
	n, udpAddr, err := udpConn.ReadFromUDP(bytes)
	if err != nil {
		return
	}

	fmt.Printf("client address: %s, received request data: %s\n", udpAddr.String(), string(bytes[0:n]))
	udpConn.WriteToUDP([]byte("welcome client"), udpAddr)
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
