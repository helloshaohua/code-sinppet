package main

import (
	"fmt"
	"sync/atomic"
)

// 序列号存储变量初始化
var code int64

// 序列号生成器
func GenerateID() int64 {
	// 序列号存储值加1
	atomic.AddInt64(&code, 1)

	// 返回序列号
	return code
}

func main() {
	// 生成10个并发序列号
	for i := 0; i < 10; i++ {
		go GenerateID()
	}

	// 输出并发序列号
	fmt.Println(GenerateID())
}
