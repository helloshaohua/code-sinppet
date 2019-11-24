// TCP Server 端设计
package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// 获取TCPAddr结构体指针对象
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8825")
	checkError(err)

	// 监听服务端口号
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// 循环接收连接请求及请求处理
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 客户端请求处理
		handleClient(tcpConn)
	}
}

// 客户端请求处理
func handleClient(conn *net.TCPConn) {
	defer conn.Close()
	buffer := make([]byte, 512)

	for {
		// 读取请求内容
		i, err := conn.Read(buffer[0:])
		if err != nil {
			return
		}

		// 获取客户端IP
		clientIP := conn.RemoteAddr()
		fmt.Printf("client address: %s, received request data: %s\n", clientIP.String(), string(buffer[0:i]))

		// 写入响应内容
		_, err = conn.Write([]byte("welcome client"))
		if err != nil {
			return
		}
	}
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
