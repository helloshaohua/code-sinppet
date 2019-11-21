### Go语言TCP协议

下面我们建立 TCP 链接来实现初步的 HTTP 协议，通过向网络主机发送 HTTP Head 请求，读取网络主机返回的信息，具体代码如下所示。

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// 命令行参数检测
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host", os.Args[0])
		os.Exit(1)
	}

	// 获取命令行参数
	service := os.Args[1]

	// 建立TCP连接
	conn, err := net.Dial("tcp", service)
	checkError(err)

	// 通过TCP连接发送请求数据
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)

	// 通过TCP连接读取响应数据
	fully, err := readFully(conn)
	checkError(err)

	// 通过标准输出，输出响应结果到命令
	fmt.Fprintf(os.Stdout, string(fully))

	// 正常退出程序
	os.Exit(0)
}

// readFully
func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()
	result := bytes.NewBuffer(nil)
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil ,err
		}
		result.Write(buf[:n])
	}
	return result.Bytes(), nil
}

// checkError 错误检测
func checkError(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "Fatal error:%s\n", e.Error())
		os.Exit(1)
	}
}
```

执行这段程序并查看执行结果：

```shell
$ go run setting-up-tcp-links.go qiniu.com:80
HTTP/1.1 200 OK
Server: nginx
Date: Thu, 21 Nov 2019 00:35:38 GMT
Content-Type: text/html
Content-Length: 612
Last-Modified: Wed, 24 Jul 2019 02:28:35 GMT
Connection: close
ETag: "5d37c253-264"
Accept-Ranges: bytes
```