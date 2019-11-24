// TCP Client端设计
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 命令行参数检测
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port\n", os.Args[0])
		os.Exit(1)
	}

	// 获取命令行参数域名或IP地址和端口号
	service := os.Args[1]

	// 获取 TCPAddr 结构体指针对象
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	// 连接到TCP
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer tcpConn.Close()

	// 获取服务器IP
	serverIP := tcpConn.RemoteAddr()

	// 发送请求数据
	n, err := tcpConn.Write([]byte("hello server"))
	checkError(err)

	// 读取响应数据
	bytes := make([]byte, 512)
	n, err = tcpConn.Read(bytes[0:])
	checkError(err)
	fmt.Fprintf(os.Stdout, "server address: %s, recevied response data: %s\n", serverIP.String(), string(bytes[0:n]))

	// 正常退出程序
	os.Exit(0)
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
