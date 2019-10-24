package main

import (
	"fmt"
	"strings"
)

func main() {
	// 字符串定义
	tracer := "死神来了，死神bye bye"

	// 搜索中文的逗号位置
	comma := strings.Index(tracer, "，")

	// 获得第二个“死神”的位置
	pos := strings.Index(tracer[comma:], "死神")

	// 输出子字符串 "死神bye bye"
	fmt.Println(comma, pos, tracer[comma+pos:])
}
