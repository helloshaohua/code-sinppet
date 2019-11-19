package main

import (
	"fmt"
	"reflect"
)

// 定义结构体
type dummy struct {
	a       int
	b       string
	float32        // 嵌入字段
	bool           // 嵌入字段
	next    *dummy // 类型为自身结构体类型
}

func main() {
	// 值包装结构体
	d := reflect.ValueOf(dummy{
		next: &dummy{a: 17995010},
	})

	// 获取字段数量
	fmt.Printf("NumField: %d\n", d.NumField())

	// 获取索引为2的字段(float32字段)
	float32Field := d.Field(2)

	// 输出字段类型
	fmt.Printf("Field Type: %s\n", float32Field.Type())

	// 根据字段名查找字段
	fmt.Printf("FieldByName(\"b\").Type: %s\n", d.FieldByName("b").Type())

	// 根据索引查找值中 next 字段的 int 字段的类型
	fmt.Printf("FieldByIndex([]int{4, 0}).Type(): %s\n", d.FieldByIndex([]int{4, 0}).Type())

	// 根据索引查找值中 next 字段的 int 字段的值
	fmt.Printf("FieldByIndex([]int{4, 0}).Type(): %d\n", d.FieldByIndex([]int{4, 0}).Int())
}
