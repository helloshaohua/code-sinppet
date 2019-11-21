### Go语言DialTCP()：网络通信

实际上，在前面 [Dial()函数](dial-func.md) 一节中介绍的 Dial() 函数其实是对 DialTCP()、DialUDP()、DialIP() 和 DialUnix() 的封装。我们也可以直接调用这些函数，它们的功能是一致的。这些函数的原型如下：

```go
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
func DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error) 
func DialIP(network string, laddr, raddr *IPAddr) (*IPConn, error)
func DialUnix(network string, laddr, raddr *UnixAddr) (*UnixConn, error)
```

在[建立TCP链接](setting-up-tcp-links.md)中基于 TCP 发送 HTTP 请求，读取服务器返回的 HTTP Head 的整个流程也可以使用下面代码所示的实现方式。

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host\n", os.Args[0])
		os.Exit(1)
	}

	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)

	bytes, err := ioutil.ReadAll(conn)
	checkError(err)

	fmt.Println(string(bytes))
	os.Exit(1)
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

```text
$ go run code/011/dial-tcp/dial-tcp.go baidu.com:80
HTTP/1.1 200 OK
Date: Thu, 21 Nov 2019 01:10:46 GMT
Server: Apache
Last-Modified: Tue, 12 Jan 2010 13:48:00 GMT
ETag: "51-47cf7e6ee8400"
Accept-Ranges: bytes
Content-Length: 81
Cache-Control: max-age=86400
Expires: Fri, 22 Nov 2019 01:10:46 GMT
Connection: Close
Content-Type: text/html
```

与之前使用 Dail() 的例子相比，这里有两个不同：

- net.ResolveTCPAddr()，用于解析地址和端口号；
- net.DialTCP()，用于建立链接。

这两个函数在 Dial() 中都得到了封装。

此外，net 包中还包含了一系列的工具函数，合理地使用这些函数可以更好地保障程序的质量。

验证 IP 地址有效性的代码如下：

```go
func ParseIP(s string) IP
```

创建子网掩码的代码如下：

```go
func IPv4Mask(a, b, c, d byte) IPMask
```

获取默认子网掩码的代码如下：

```go
func (ip IP) DefaultMask() IPMask
```

根据域名查找 IP 的代码如下：

```go
func ResolveIPAddr(network, address string) (*IPAddr, error)
func LookupHost(host string) (addrs []string, err error)
```

ResolveIPAddr:

```go
addr, _ := net.ResolveIPAddr("ip", "baidu.com")
// 39.156.69.79
fmt.Printf("%+v\n", addr)
```

LookupHost:

```go
addr, _ := net.LookupHost("baidu.com")
// [39.156.69.79 220.181.38.148]
fmt.Printf("%+v\n", addr)
```
