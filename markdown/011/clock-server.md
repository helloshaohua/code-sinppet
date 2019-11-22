### 示例：并发时钟服务器

网络是一个自然使用并发的领域，因为服务器通常一次处理很多来自客户端的连接，每一个客户端通常和其他客户端保持独立。本节介绍 net 包，它提供构建客户端和服务器程序的组件，这些程序通过 TCP、UDP 或者 UNIX 套接字进行通信。net/http 包就是在 net 包基础上构建的。

#### 顺序时钟服务器

所谓的有顺序时针服务器就是它以每秒钟一次的频率向客户端发送当前时间，每次只接收一个客户端请求，当这个客户端请求完成，断开连接了，下一个客户端连接才被允许连接，代码如下所示：

```go
//clock1 是一个定期报告时间的 TCP 服务器
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// 一次处理一个连接
		handlerConn(conn)
	}
}

// 处理连接
func handlerConn(conn net.Conn) {
	defer conn.Close()

	for {
		_, err := io.WriteString(conn, time.Now().Format("15:03:04\n"))
		if err != nil {
			// 例如：连接断开
			return
		}

		time.Sleep(1 * time.Second)
	}
}
```

Listen 函数创建一个 net.Listener 对象，它在一个网络端口上监听进来的连接，这里是 TCP 端口 localhost:8000。监听器的 Accept 方法被阻塞，直到有连接请求进来，然后返回 net.Conn 对象来代表一个连接。

handleConn 函数处理一个完整的客户连接。在循环里，它将 time.Now() 获取的当前时间发送给客户端。因为 net.Conn 满足 io.Writer 接口，所以可以直接向它进行写入。

当写入失败时循环结束，很多时候是客户端断开连接，这时 handleconn 函数使用延迟的 Close 调用关闭自己这边的连接，然后继续等待下一个连接请求。

time.Time.Format 方法提供了格式化日期和时间信息的方式。它的参数是一个模板，指示如何格式化一个参考时间，具体如 Mon Jan 2 03:04:05PM 2006 UTC-0700 这样的形式。参考时间有 8 个部分（本周第几天、月、本月第几天，等等）。

它们可以以任意的组合和对应数目的格式化字符出现在格式化模板中，所选择的日期和时间将通过所选择的格式进行显示。这里只使用时间的小时、分钟和秒部分。time 包定义了许多标准时间格式的模板，如 time.RFC1123。相反，当解析一个代表时间的字符串的时候使用相同的机制。

为了连接到服务器，需要一个像 nc("netcat") 这样的程序，或者一个用来操作网络连接的标准工具，下面使用 nc 连接到 TCP 服务：

```shell
$ go build clock1.go
$ ./clock1 &
$ nc localhost 8080
19:03:13
19:03:14
19:03:15
19:03:16
19:03:17
19:03:18
```

客户端显示每秒从服务器发送的时间，如果这个时候就开启另一个终端再次执行 `nc localhost 8080` 你会发现它一直处于等待状态，直到使用 Control+C 快捷键中断另一个终端，UNIX 系统 shell 上面回显为 ^C，新开启的终端才会输出时间。

如果系统上没有安装 nc 或 netcat，可以使用 telnet 或者一个使用 net.Dial 实现的 Go 版的 netcat 来连接 TCP 服务器：

```go
package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(conn)
}

func mustCopy(conn net.Conn) {
	if _, err := io.Copy(os.Stdout, conn); err != nil {
		log.Fatal(err)
	}
}
```

这个程序从网络连接中读取，然后写到标准输出，直到到达 EOF 或者岀错。在不同的终端上同时运行两个客户端，一个显示在左边，一个在右边：

```shell
$ go build netcat.go
```

![顺序时钟服务器](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191122_3.png)

使用killall结束进程。

```shell
$ killall clock1
```

killall 命令是 UNIX 的一个实用程序，用来终止所有指定名字的进程。


第二个客户端必须等到第一个结束才能正常工作，因为服务器是顺序的，一次只能处理一个客户请求。让服务器支持并发只需要一个很小的改变，在调用 handlerConn 的地方添加一个 go 关键字，使它在自己的 goroutine 内执行。

#### 并发时钟服务器

```go
//clock2 是一个并发的定期报告时间的 TCP 服务器
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// 并发处理连接
		go handlerConn(conn)
	}
}

// 处理连接
func handlerConn(conn net.Conn) {
	defer conn.Close()

	for {
		_, err := io.WriteString(conn, time.Now().Format("15:04:05\n"))
		if err != nil {
			// 例如：连接断开
			return
		}

		time.Sleep(1 * time.Second)
	}
}
```

现在，多个客户端可以同时接收到时间：

```shell
$ go build clock2.go
$ ./clock2 &
```

![并发时钟服务器](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191122_4.gif)



