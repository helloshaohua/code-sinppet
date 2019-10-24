package main

import "fmt"

func main() {
	// 字符串变量定义
	angel := "Heros never die"

	// 将字符串转为字符串切片
	angelBytes := []byte(angel)

	// 利用循环，将 never 单词替换为空格
	for i := 5; i <= 10; i++ {
		angelBytes[i] = ' '
	}

	fmt.Println(string(angelBytes))
}
