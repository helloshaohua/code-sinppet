package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 声明 dog 结构体
	type dog struct {
		LegCount int
	}

	// 初始化结构体实例得到一个结构体指针
	d := &dog{}

	// 获取 dog 实例的反射值对象
	valueOfDog := reflect.ValueOf(d)

	// 获取 dog 实例地址的元素
	valueOfDog = valueOfDog.Elem()

	// 获取 LegCount 字段的值
	vLegCount := valueOfDog.FieldByName("LegCount")

	// 修改或设置 LegCount 字段的值
	vLegCount.SetInt(1024)

	// 输出修改后的值
	fmt.Println(vLegCount.Int())

	// 输出结构体值
	fmt.Printf("%+v\n", d)
}
