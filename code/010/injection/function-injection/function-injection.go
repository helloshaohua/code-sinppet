package main

import (
	"fmt"

	"github.com/codegangsta/inject"
)

type S1 interface {
}

type S2 interface {
}

func Format(name string, company S1, level S2, age int) {
	fmt.Printf("name: %s, company: %s, level: %d, age: %d\n", name, company, level, age)
}

func main() {
	// 控制实例的创建
	i := inject.New()

	// 实参注入
	i.Map("张三")
	i.MapTo("阿里巴巴", (*S1)(nil))
	i.MapTo(8, (*S2)(nil))
	i.Map(18)

	// 函数反转调用
	i.Invoke(Format)
}
