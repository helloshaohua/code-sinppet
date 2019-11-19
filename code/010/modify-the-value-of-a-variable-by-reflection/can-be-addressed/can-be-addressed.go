package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 声明整型变量 a 并赋初值
	var a int = 1024

	// 获取变量 a 的反射值对象(a的地址)
	valueOfA := reflect.ValueOf(&a)

	// 取出 a 地址的元素(a的值)
	valueOfA = valueOfA.Elem()

	// 修改 a 的值为 2048
	valueOfA.SetInt(2048)

	// 打印 a 的值
	fmt.Println(valueOfA.Interface())
}
