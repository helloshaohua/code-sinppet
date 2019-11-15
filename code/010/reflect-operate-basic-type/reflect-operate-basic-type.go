package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 修改字符串值
	username := "张三"
	reflect.ValueOf(&username).Elem().SetString("李四")
	fmt.Println(username)

	// 修改整型值
	age := 22
	reflect.ValueOf(&age).Elem().SetInt(18)
	fmt.Println(age)

	// 修改浮点数值
	weight := 65.0
	reflect.ValueOf(&weight).Elem().SetFloat(60.83)
	fmt.Println(weight)

	// 修改布尔值
	died := false
	reflect.ValueOf(&died).Elem().SetBool(true)
	fmt.Println(died)

	// 修改字节值
	say := []byte("你好")
	reflect.ValueOf(&say).Elem().SetBytes([]byte("Hello World, 你好！"))
	fmt.Println(string(say))
}
