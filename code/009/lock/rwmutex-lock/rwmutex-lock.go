package main

import (
	"fmt"
	"sync"
)

// 变量列表声明
var (
	count      int          // 逻辑中使用的计数变量
	countGuard sync.RWMutex // 与变量对应的读写互斥锁使用
)

// GetCount 获取计数器
func GetCount() int {
	// 使用读写锁锁定
	countGuard.RLock()

	// 在函数退出时解除锁定
	defer countGuard.RUnlock()

	// 返回计数
	return count
}

// SetCount 设置计数器
func SetCount(number int) {
	// 改变计数器值
	count = number
}

func main() {
	// 可以进行并发安全的设置计数器值
	SetCount(1)

	// 可以进行并发安全的获取计数器值
	fmt.Println(GetCount())
}
