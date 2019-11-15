package main

import (
	"fmt"
	"reflect"
)

// Student 学生结构体类型
type Student struct {
	Name   string  // 姓名
	Age    int     // 年龄
	Sex    uint8   // 1 男， 2 女
	Weight float64 // 体重
}

func main() {
	// 用户实例
	student := Student{Name: "张三", Age: 18, Sex: 1, Weight: 68.57}

	// 获取实例的值
	v := reflect.ValueOf(&student).Elem()

	// 获取实例的值的类型
	t := v.Type()

	// 通过【实例的值】的方法【NumField】获取字段的个数
	for i := 0; i < v.NumField(); i++ {
		// 通过【实例的值】的方法【Field】获取字段
		field := v.Field(i)

		// 字段名称
		fieldName := t.Field(i).Name

		// 字段类型
		fieldType := field.Type()

		// 字段值
		fieldValue := field.Interface()

		// 输出结构体成员信息
		fmt.Printf("字段索引值：%d, 字段名称：%s, 字段类型：%s, 字段值：%v\n", i, fieldName, fieldType, fieldValue)
	}
}
