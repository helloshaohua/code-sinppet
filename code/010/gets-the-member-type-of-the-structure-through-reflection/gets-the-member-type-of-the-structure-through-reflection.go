package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 声明【猫】结构体
	type cat struct {
		Name string
		Type int `json:"type" id:"100"` // 带有标签字段
	}

	// 创建结构体实例
	instance := cat{Name: "mini", Type: 1}

	// 获取结构体实例的反射实例的反射类型对象
	typeOfCat := reflect.TypeOf(instance)

	// 遍历结构体所有成员
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取结构体字段类型
		field := typeOfCat.Field(i)

		// 输出成员名和标签
		fmt.Printf("name: %s, tag: '%s'\n", field.Name, field.Tag)
	}

	// 通过字段名，获取字段类型信息
	if field, exist := typeOfCat.FieldByName("Type"); exist {
		// 从标签中取出需要的具体标签
		fmt.Printf("json tag statement: %s, id tag statement: %s\n", field.Tag.Get("json"), field.Tag.Get("id"))
	}
}
