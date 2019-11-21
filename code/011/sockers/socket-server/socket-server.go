package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	service = ":8299"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	request := make([]byte, 128)
	for {
		readLen, err := conn.Read(request)
		if err != nil {
			fmt.Println(err)
			break
		}
		if readLen == 0 {
			break // 客户端已关闭连接
		} else if strings.TrimSpace(string(request[:readLen])) == "timestamp" {
			daytime := strconv.FormatInt(time.Now().Unix(), 10)
			conn.Write([]byte(daytime))
		} else {
			daytime := time.Now().String()
			conn.Write([]byte(daytime))
		}
		request = make([]byte, 128) // 清除上次读取的内容
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
