package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 声明一个空结构体
	type cat struct{}

	// 创建空结构体实例
	instance := &cat{}

	// 获取结构体实例反射类型对象
	typeOfCat := reflect.TypeOf(instance)

	// 显示反射类型对象的名称和种类
	fmt.Printf("name: '%v', kind: '%v'\n", typeOfCat.Name(), typeOfCat.Kind())

	// 获取类型对象的元素
	typeOfCat = typeOfCat.Elem()

	// 显示反射类型对象的名称和种类
	fmt.Printf("name: '%v', kind: '%v'\n", typeOfCat.Name(), typeOfCat.Kind())
}
