package main

import (
	"fmt"
	"strings"
)

// 字符串处理函数，传入字符串切片和处理链
func StringProcess(list []string, chain []func(string) string) []string {
	for index, str := range list {
		result := str
		for _, function := range chain {
			result = function(result)
		}
		list[index] = result
	}
	return list
}

func main() {
	// 待处理的字符串列表
	list := []string{
		"go scanner",
		"go parser",
		"go compiler",
		"go printer",
		"go formater",
	}

	// 函数处理链
	chains := []func(string) string{
		removeStringPrefix,
		strings.TrimSpace,
		strings.ToUpper,
	}

	// 处理字符串
	list = StringProcess(list, chains)

	// 输出处理好的字符串
	for _, str := range list {
		fmt.Println(str)
	}
}

// 自定义的移除前缀的处理函数
func removeStringPrefix(in string) (out string) {
	return strings.TrimPrefix(in, "go")
}
