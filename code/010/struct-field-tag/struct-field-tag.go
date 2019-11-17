package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 声明【猫】结构体
	type cat struct {
		Name string
		Type int `json:"type" id:"100"`
	}

	// 返回结构体类型值反射对象
	typeOfCat := reflect.TypeOf(cat{})

	// 获取对应的结构体字段的信息
	if field, exist := typeOfCat.FieldByName("Type"); exist {
		fmt.Printf("Structure field Type marks the json value: %s\n", field.Tag.Get("json"))
	}
}
