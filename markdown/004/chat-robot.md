### 示例：聊天机器人

结合咱们之前的学习，本节带领大家来编写一个聊天机器人的雏形，下面的代码中展示了一个简单的聊天程序。

```go
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
```

这个聊天程序在问候用户之后会不断地询问“是否可以帮忙”，但是实际上它什么忙也帮不上，因为它现在什么也听不懂，除了 nothing 和 bye，一看到这两个词，它就会与用户“道别”，停止运行，现在试运行一下这个命令源码文件：

```shell
$ go run test.go
-> Please input your name: 
wu.shaohua@foxmail.com
Hello, wu.shaohua@foxmail.com! What can I do for you?
-> hello world
Haha, HELLO WORLD
->bye
Bye!
```

注意，其中的“->”符号之后的内容是我们输入的。

