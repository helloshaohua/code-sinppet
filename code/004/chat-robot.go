package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const NewLineCharacterOccupyLen = 1

func main() {
	// 准备从标准输入读取数据
	reader := bufio.NewReader(os.Stdin)

	// 提示用户输入用户名
	fmt.Println("Please input your name: ")

	// 读取数据直到碰到 \n 为止
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		os.Exit(0)
	} else {
		// 用切片操作删除最后面的 \n
		username = username[:len(username)-NewLineCharacterOccupyLen]
		fmt.Printf("Hello, %s! What can I do for you?\n", username)
	}

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("An error occurred: %s\n", err)
			continue
		}

		// 获取输入内容
		s = s[:len(s)-NewLineCharacterOccupyLen]
		switch s {
		case "":
			continue
		case "bye", "nothing":
			fmt.Printf("Bye!\n")
			os.Exit(0)
		default:
			fmt.Printf("Haha, %s\n", strings.ToUpper(s))
		}
	}
}
