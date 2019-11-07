package main

import (
	"code-snippet/code/008/clsfactory/base"
	_ "code-snippet/code/008/clsfactory/cls1"
	_ "code-snippet/code/008/clsfactory/cls2"
)

func main() {
	// 根据字符串动态创建一个 Class1 实例
	c1 := base.Create("Class1")
	c1.Do()

	// 根据字符串动态创建一个 Class2 实例
	c2 := base.Create("Class2")
	c2.Do()
}
