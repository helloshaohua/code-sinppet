package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func main() {
	// 创建一个无缓冲字符串通道
	ch := make(chan string)

	// 并发执行服务器逻辑
	go RPCServer(ch)

	// 客户端请求数据和接收数据
	response, err := RPCClient(ch, "hi")
	if err != nil {
		log.Printf("RPCClient error: %s\n", err)
	} else {
		fmt.Printf("client received: %s\n", response)
	}
}

// 模拟RPC服务器端接收客户端请求和响应
func RPCServer(ch chan string) {
	for {
		// 接收客户端请求
		data := <-ch

		// 打印接收到的数据
		fmt.Printf("Server received: %s\n", data)

		// 通过睡眠函数让程序执行阻塞2秒的任务
		time.Sleep(2 * time.Second)

		// 反馈给客户端收到了
		ch <- "ok"
	}
}

// 模拟RPC客户端请求和接收消息
func RPCClient(ch chan string, data string) (string, error) {
	// 向服务器发送请求
	ch <- data

	// 等待服务器返回
	select {
	case response := <-ch: // 接收通道响应数据
		return response, nil
	case <-time.After(time.Second): // 超时
		return "", errors.New("timeout")
	}
}
