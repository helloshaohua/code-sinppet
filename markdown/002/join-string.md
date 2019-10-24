### Go语言字符串拼接（连接）

连接字符串这么简单，还需要学吗？确实，Go 语言和大多数其他语言一样，使用+对字符串进行连接操作，非常直观。

但问题来了，好的事物并非完美，简单的东西未必高效。除了加号连接字符串，Go 语言中也有类似于 StringBuilder 的机制来进行高效的字符串连接，例如：

```go
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
```

bytes.Buffer 是可以缓冲并可以往里面写入各种字节数组的。字符串也是一种字节数组，使用 WriteString() 方法进行写入。

将需要连接的字符串，通过调用 WriteString() 方法，写入 stringBuilder 中，然后再通过 stringBuilder.String() 方法将缓冲转换为字符串。


