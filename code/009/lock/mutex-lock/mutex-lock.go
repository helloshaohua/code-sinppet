package main

import (
	"fmt"
	"sync"
)

// 变量列表声明
var (
	count      int        // 逻辑中使用的计数变量
	countGuard sync.Mutex // 与变量对应的互斥锁使用
)

// GetCount 获取计数器
func GetCount() int {
	// 锁定
	countGuard.Lock()

	// 在函数退出时解除锁定
	defer countGuard.Unlock()

	// 返回计数
	return count
}

// SetCount 设置计数器
func SetCount(number int) {
	// 锁定
	countGuard.Lock()

	// 改变计数器值
	count = number

	// 解除锁定
	countGuard.Unlock()
}

func main() {
	// 可以进行并发安全的设置计数器值
	SetCount(1)

	// 可以进行并发安全的获取计数器值
	fmt.Println(GetCount())
}
