package main

import (
	"fmt"

	"github.com/go-macaron/inject"
)

type S1 interface {
}

type S2 interface {
}

type Staff struct {
	Name    string `inject:"name"`
	Company S1     `inject:"company"`
	Level   S2     `inject:"level"`
	Age     int    `inject:"age"`
}

func main() {
	// 创建被注入实例
	s := Staff{}

	// 控制实例的创建
	si := inject.New()

	// 初始化注入值
	si.Map("张三")
	si.MapTo("阿里巴巴", (*S1)(nil))
	si.MapTo(8, (*S2)(nil))
	si.Map(18)

	// 实现对 struct 注入
	si.Apply(&s)

	// 打印注入后结构体
	fmt.Printf("s: %+v\n", s)
}
