### 创建一个HTTP请求

Go语言提供的 http 包里也大量使用了类型方法，Go语言使用 http 包进行 HTTP 的请求，使用 http 包的 NewRequest() 方法可以创建一个 HTTP 请求，填充请求中的 http 头（req.Header），再调用 http.Client 的 Do 方法，将传入的 HTTP 请求发送出去。

下面代码演示创建一个 HTTP 请求，并且设定 HTTP 头。
```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	client := new(http.Client)

	// 创建一个http请求
	request, err := http.NewRequest("POST", "http://www.163.com/", strings.NewReader("key=value"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 为请求头添加信息
	request.Header.Add("User-Agent", "myClient")

	// 开始请求
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 关闭响应体
	defer response.Body.Close()

	// 读取响应体数据
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 输出响应体
	fmt.Println(string(bytes))
}
```

代码执行结果如下：
```text
<html>
<head><title>405 Not Allowed</title></head>
<body bgcolor="white">
<center><h1>405 Not Allowed</h1></center>
<hr><center>nginx</center>
</body>
</html>
```

代码说明如下：

- 第 12 行，实例化 HTTP 的客户端，请求需要通过这个客户端实例发送。
- 第 15 行，使用 POST 方式向网易的服务器创建一个 HTTP 请求，第三个参数为 HTTP 的 Body 部分，Body 部分的内容来自字符串，但参数只能接受 io.Reader 类型，因此使用 strings.NewReader() 创建一个字符串的读取器，返回的 io.Reader 接口作为 http 的 Body 部分供 NewRequest() 函数读取，创建请求只是构造一个请求对象，不会连接网络。
- 第 23 行，为创建好的 HTTP 请求的头部添加 User-Agent，作用是表明用户的代理特性。
- 第 26 行，使用客户端处理请求，此时 client 将 HTTP 请求发送到网易服务器，服务器响应请求后，将信息返回并保存到 resp 变量中。
- 第 37 ~ 45 行，读取响应的 Body 部分并打印。

由于我们构造的请求不是网易服务器所支持的类型，所以服务器返回操作不被运行的 405 错误。

在本例子第 23 行中使用的 req.Header 的类型为 http.Header，就是典型的自定义类型，并且拥有自己的方法，http.Header 的部分定义如下：

```go
type Header map[string][]string
func (h Header) Add(key, value string) {
    textproto.MIMEHeader(h).Add(key, value)
}
func (h Header) Set(key, value string) {
    textproto.MIMEHeader(h).Set(key, value)
}
func (h Header) Get(key string) string {
    return textproto.MIMEHeader(h).Get(key)
}
```

代码说明如下：
- 第 1 行，Header 实际是一个以字符串为键、字符串切片为值的映射。
- 第 3 行，Add() 为 Header 的方法，map 是一个引用类型，因此即便使用 (h Header) 的非指针接收器，也可以修改 map 的值。

为类型添加方法的过程是一个语言层特性，使用类型方法的代码经过编译器编译后的代码运行效率与传统的面向过程或面向对象的代码没有任何区别，因此，为了代码便于理解，可以在编码时使用Go语言的类型方法特性。

