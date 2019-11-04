package main

import "fmt"

// 声明一个结构体
type class struct {
}

// 给结构体添加Do方法
func (c *class) Do(v int) {
	fmt.Println("call method do:", v)
}

// 普通函数的 Do
func FuncDo(v int) {
	fmt.Println("call function do:", v)
}

func main() {
	// 声明一个函数回调
	var delegate func(int)

	// 创建结构体实例返回结构体指针
	c := new(class)

	// 将回调设为c的 Do方法
	delegate = c.Do

	// 调用
	delegate(100)

	// 将回调设为普通函数
	delegate = FuncDo

	// 调用
	delegate(100)
}
