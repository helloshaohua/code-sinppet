package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// 服务逻辑，传入地址和退出的通道
func server(address string, exit chan int) {
	// 根据给定的地址进行侦听
	listen, err := net.Listen("tcp", address)

	// 如果侦听发生错误，打印错误并退出
	if err != nil {
		fmt.Println(err.Error())
		exit <- 1
	}

	// 打印侦听地址，表示侦听成功
	fmt.Printf("listen: %s\n", address)

	// 延迟关闭侦听器
	defer listen.Close()

	// 侦听循环
	for {
		// 新连接没有到来时，Accept是阻塞的
		conn, err := listen.Accept()

		// 发生任何接受错误，打印错误并忽略本次连接
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// 根据连接开启会话，这个过程需要并行执行
		go handlerSession(conn, exit)
	}
}

// 连接的会话逻辑
func handlerSession(conn net.Conn, exit chan int) {
	// 开始会话处理提示信息
	fmt.Println("Session started:")

	// 创建一个网络连接数据的读取器
	reader := bufio.NewReader(conn)

	// 接收数据的循环
	for {
		// 读取字符串，直接碰到回车
		data, err := reader.ReadString('\n')
		if err != nil {
			// 发生错误
			fmt.Println("Session closed")
			conn.Close()
			break
		} else {
			// 去掉字符串尾部的回车符
			data = strings.TrimSpace(data)

			// 处理Telnet指令
			if !processTelnetCommand(data, exit) {
				conn.Close()
				break
			}

			// 回声逻辑，发什么数据，原样返回
			conn.Write([]byte(data + "\r\n"))
		}
	}
}

// 命令处理
func processTelnetCommand(data string, exit chan int) bool {
	// @close指令表示终止本次会话
	if strings.HasPrefix(data, "@close") {
		// 提示终止本次会话
		fmt.Println("Session closed")

		// 告诉外部需要断开连接
		return false
	}

	// @shutdown指令表示终止服务进程
	if strings.HasPrefix(data, "@shutdown") {
		// 提示终止服务进程
		fmt.Println("Server shutdown")

		// 向通道中写入0，main函数中的阻塞等待接收方会处理
		exit <- 0

		// 告诉外部需要断开连接
		return false
	}

	// 打印用户输入的字符串
	fmt.Println(data)

	// 告诉外部不需要断开
	return true
}

func main() {
	// 创建一个程序结束状态的通道
	exit := make(chan int)

	// 将服务器并发运行
	go server("127.0.0.1:8001", exit)

	// 通道阻塞，等待接收返回值
	code := <-exit

	// 标记程序返回值并退出
	os.Exit(code)
}
