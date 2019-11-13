package main

import (
	"fmt"
	"time"
)

func main() {
	// 声明一个退出用的通道
	exit := make(chan bool)

	// 打印开始
	fmt.Println("start")

	// 过1秒钟后，调用匿名函数
	time.AfterFunc(time.Second, func() {
		// 1秒钟后，打印结束
		fmt.Println("one second after-func")

		// 通知main函数的 goroutine 已经结束
		exit <- true
	})

	// 等待结束
	<-exit
}
