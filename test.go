package main

import (
	"fmt"

	"github.com/go-macaron/inject"
)

type S1 interface {
}

type S2 interface {
}

func Format(name string, company S1, level S2, age int) {
	fmt.Printf("name: %s, company: %s, level: %s, age: %d\n", name, company, level, age)
}

func main() {
	// 控制实例的创建
	i := inject.New()
	// 实参注入
	i.Map("tom")
	i.MapTo("tencent", (*S1)(nil))
	i.MapTo("T4", (*S2)(nil))
	i.Map(18)
	fmt.Printf("%+v\n", (*S1)(nil))
	// 函数反转调用
	values, err := i.Invoke(Format)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(values)
	}
}
