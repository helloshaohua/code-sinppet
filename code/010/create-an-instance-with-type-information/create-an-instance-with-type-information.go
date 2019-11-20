package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 声明整型变量 a
	var a int

	// 获取变量 a 的反射类型对象
	typeOfA := reflect.TypeOf(a)

	// 根据反射类型对象创建类型实例
	instance := reflect.New(typeOfA)

	// 输出Value的类型和种类
	fmt.Println(instance.Type(), instance.Kind())

	// 修改类型实例值
	instance.Elem().SetInt(2048)

	// 输出类型值
	fmt.Printf("instance: %+v\n", instance.Elem().Int())
}
