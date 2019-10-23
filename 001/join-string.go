package main

import (
	"bytes"
	"fmt"
)

func main() {
	hammer := "吃我一锤"
	sickle := "死吧"
	// 声明字节缓冲
	var stringBuilder bytes.Buffer
	// 把字符串写入缓冲
	stringBuilder.WriteString(hammer)
	stringBuilder.WriteString(sickle)
	// 将缓冲以字符串形式输出
	fmt.Println(stringBuilder.String())
}
