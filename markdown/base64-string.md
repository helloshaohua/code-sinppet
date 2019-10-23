### Go语言Base64编码——电子邮件的基础编码格式

Base64 编码是常见的对 8 比特字节码的编码方式之一。Base64 可以使用 64 个可打印字符来表示二进制数据，电子邮件就是使用这种编码。

Go 语言的标准库自带了 Base64 编码算法，通过几行代码就可以对数据进行编码，示例代码如下。

```go
package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	// 需求处理的字符串
	message := "Away from keyboard. https://golang.org/"

	// 编码消息
	encodeMessage := base64.StdEncoding.EncodeToString([]byte(message))

	// 输出编码完成的消息
	fmt.Println(encodeMessage)

	// 解码消息
	data, err := base64.StdEncoding.DecodeString(encodeMessage)
	if err != nil {
		log.Fatalf("base64.StdEncoding.DecodeString error:%s\n", err)
	}

	// 输出解码完成的消息
	fmt.Println(string(data))
}
```

程序输出如下：

```go
package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	// 需求处理的字符串
	message := "Away from keyboard. https://golang.org/"

	// 编码消息
	encodeMessage := base64.StdEncoding.EncodeToString([]byte(message))

	// 输出编码完成的消息
	fmt.Println(encodeMessage)

	// 解码消息
	data, err := base64.StdEncoding.DecodeString(encodeMessage)
	if err != nil {
		log.Fatalf("base64.StdEncoding.DecodeString error:%s\n", err)
	}

	// 输出解码完成的消息
	fmt.Println(string(data))
}
```

代码说明如下：
- 第 11 行为需要编码的消息，消息可以是字符串，也可以是二进制数据。
- 第 14 行，base64 包有多种编码方法，这里使用 base64.StdEnoding 的标准编码方法进行编码。传入的字符串需要转换为字节数组才能供这个函数使用。
- 第 17 行，编码完成后一定会输出字符串类型，打印输出。
- 第 20 行，解码时可能会发生错误，使用 err 变量接收错误。
- 第 22 行，出错时，打印错误。
- 第 26 行，正确时，将返回的字节数组（[]byte）转换为字符串。
