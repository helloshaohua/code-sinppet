### Go语言模拟远程过程调用

服务器开发中会使用RPC（Remote Procedure Call，远程过程调用）简化进程间通信的过程。RPC 能有效地封装通信过程，让远程的数据收发通信过程看起来就像本地的函数调用一样。

本例中，使用通道代替 Socket 实现 RPC 的过程。客户端与服务器运行在同一个进程，服务器和客户端在两个 goroutine 中运行。

我们先给出完整代码，然后再详细分析每一个部分。

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func main() {
	// 创建一个无缓冲字符串通道
	ch := make(chan string)

	// 并发执行服务器逻辑
	go RPCServer(ch)

	// 客户端请求数据和接收数据
	response, err := RPCClient(ch, "hi")
	if err != nil {
		log.Printf("RPCClient error: %s\n", err)
	} else {
		fmt.Printf("client received: %s\n", response)
	}
}

// 模拟RPC服务器端接收客户端请求和响应
func RPCServer(ch chan string) {
	for {
		// 接收客户端请求
		data := <-ch

		// 打印接收到的数据
		fmt.Printf("Server received: %s\n", data)

		// 反馈给客户端收到了
		ch <- "ok"
	}
}

// 模拟RPC客户端请求和接收消息
func RPCClient(ch chan string, data string) (string, error) {
	// 向服务器发送请求
	ch <- data

	// 等待服务器返回
	select {
	case response := <-ch: // 接收通道响应数据
		return response, nil
	case <-time.After(time.Second): // 超时
		return "", errors.New("timeout")
	}
}
```

#### 客户端请求和接收响应封装

下面的代码封装了向服务器请求数据，等待服务器返回数据，如果请求方超时，该函数还会处理超时逻辑。

模拟 RPC 的代码：

```go
// 模拟RPC客户端请求和接收消息
func RPCClient(ch chan string, data string) (string, error) {
	// 向服务器发送请求
	ch <- data

	// 等待服务器返回
	select {
	case response := <-ch: // 接收通道响应数据
		return response, nil
	case <-time.After(time.Second): // 超时
		return "", errors.New("timeout")
	}
}
```

代码说明如下：
- 第 4 行，模拟 socket 向服务器发送一个字符串信息。服务器接收后，结束阻塞执行下一行。
- 第 7 行，使用 select 开始做多路复用。注意，select 虽然在写法上和 switch 一样，都可以拥有 case 和 default。但是 select 关键字后面不接任何语句，而是将要复用的多个通道语句写在每一个 case 上，如第 9 行和第 11 行所示。
- 第 10 行，使用了 time 包提供的函数 After()，从字面意思看就是多少时间之后，其参数是 time 包的一个常量，time.Second 表示 1 秒。time.After 返回一个通道，这个通道在指定时间后，通过通道返回当前时间。
- 第 11 行，在超时时，返回超时错误。

RPCClient() 函数中，执行到 select 语句时，第 9 行和第 11 行的通道操作会同时开启。如果第 9 行的通道先返回，则执行第 10 行逻辑，表示正常接收到服务器数据；如果第 11 行的通道先返回，则执行第 12 行的逻辑，表示请求超时，返回错误。

#### 服务器接收和反馈数据封装

服务器接收到客户端的任意数据后，先打印再通过通道返回给客户端一个固定字符串，表示服务器已经收到请求。

```go
// 模拟RPC服务器端接收客户端请求和响应
func RPCServer(ch chan string) {
	for {
		// 接收客户端请求
		data := <-ch

		// 打印接收到的数据
		fmt.Printf("Server received: %s\n", data)

		// 反馈给客户端收到了
		ch <- "ok"
	}
}
```

代码说明如下：
- 第 3 行，构造出一个无限循环。服务器处理完客户端请求后，通过无限循环继续处理下一个客户端请求。
- 第 5 行，通过字符串通道接收一个客户端的请求。
- 第 8 行，将接收到的数据打印出来。
- 第 11 行，给客户端反馈一个字符串。

运行整个程序，客户端可以正确收到服务器返回的数据，客户端 RPCClient() 函数的代码按下面代码中加粗部分的分支执行。

```go
// 等待服务器返回
select {
case response := <-ch: // 接收通道响应数据
    return response, nil
case <-time.After(time.Second): // 超时
    return "", errors.New("timeout")
}
```

程序输出如下：

```text
Server received: hi
client received: ok
```

#### 模拟超时

上面的例子虽然有客户端超时处理，但是永远不会触发，因为服务器的处理速度很快，也没有真正的网络延时或者“服务器宕机”的情况。因此，为了展示 select 中超时的处理，在服务器逻辑中增加一条语句，故意让服务器延时处理一段时间，造成客户端请求超时，代码如下：

```go
// 模拟RPC服务器端接收客户端请求和响应
func RPCServer(ch chan string) {
	for {
		// 接收客户端请求
		data := <-ch

		// 打印接收到的数据
		fmt.Printf("Server received: %s\n", data)

		// 通过睡眠函数让程序执行阻塞2秒的任务
		time.Sleep(2 * time.Second)

		// 反馈给客户端收到了
		ch <- "ok"
	}
}
```

第 11 行中，time.Sleep() 函数会让 goroutine 执行暂停 2 秒。使用这种方法模拟服务器延时，造成客户端超时。客户端处理超时 1 秒时通道就会返回：

```go
// 等待服务器返回
select {
case response := <-ch: // 接收通道响应数据
    return response, nil
case <-time.After(time.Second): // 超时
    return "", errors.New("timeout")
}
```

上面代码中，第5行中的代码就会被执行。


#### 主流程

主流程中会创建一个无缓冲的字符串格式通道。将通道传给服务器的 RPCServer() 函数，这个函数并发执行。使用 RPCClient() 函数通过 ch 对服务器发出 RPC 请求，同时接收服务器反馈数据或者等待超时。参考下面代码：

```go
func main() {
	// 创建一个无缓冲字符串通道
	ch := make(chan string)

	// 并发执行服务器逻辑
	go RPCServer(ch)

	// 客户端请求数据和接收数据
	response, err := RPCClient(ch, "hi")
	if err != nil {
		log.Printf("RPCClient error: %s\n", err)
	} else {
		fmt.Printf("client received: %s\n", response)
	}
}
```

代码说明如下：
- 第 3 行，创建无缓冲的字符串通道，这个通道用于模拟网络和 socke t概念，既可以从通道接收数据，也可以发送。
- 第 6 行，并发执行服务器逻辑。服务器一般都是独立进程的，这里使用并发将服务器和客户端逻辑同时在一个进程内运行。
- 第 9 行，使用 RPCClient() 函数，发送“hi”给服务器，同步等待服务器返回。
- 第 11 行，如果通信过程发生错误，打印错误。
- 第 13 行，正常接收时，打印收到的数据。

