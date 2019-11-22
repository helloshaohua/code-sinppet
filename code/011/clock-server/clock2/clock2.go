//clock2 是一个并发的定期报告时间的 TCP 服务器
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// 并发处理连接
		go handlerConn(conn)
	}
}

// 处理连接
func handlerConn(conn net.Conn) {
	defer conn.Close()

	for {
		_, err := io.WriteString(conn, time.Now().Format("15:04:05\n"))
		if err != nil {
			// 例如：连接断开
			return
		}

		time.Sleep(1 * time.Second)
	}
}
