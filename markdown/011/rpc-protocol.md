### Go语言RPC协议：远程过程调用

Go语言中 RPC（Remote Procedure Call，远程过程调用）是一种通过网络从远程计算机程序上请求服务，而不需要了解底层网络细节的应用程序通信协议。RPC 协议构建于 TCP 或 UDP，或者是 HTTP 之上，允许开发者直接调用另一台计算机上的程序，而开发者无需额外地为这个调用过程编写网络通信相关代码，使得开发包括网络分布式程序在内的应用程序更加容易。

RPC 采用客户端—服务器（Client/Server）的工作模式。请求程序就是一个客户端（Client），而服务提供程序就是一个服务器（Server）。当执行一个远程过程调用时，客户端程序首先发送一个带有参数的调用信息到服务端，然后等待服务端响应。

在服务端，服务进程保持睡眠状态直到客户端的调用信息到达为止。当一个调用信息到达时，服务端获得进程参数，计算出结果，并向客户端发送应答信息，然后等待下一个调用。最后，客户端接收来自服务端的应答信息，获得进程结果，然后调用执行并继续进行。

在 Go 中，标准库提供的 net/rpc 包实现了 RPC 协议需要的相关细节，开发者可以很方便地使用该包编写 RPC 的服务端和客户端程序，这使得用 Go语言开发的多个进程之间的通信变得非常简单。

net/rpc 包允许 RPC 客户端程序通过网络或是其他 I/O 连接调用一个远端对象的公开方法（必须是大写字母开头、可外部调用的）。在 RPC 服务端，可将一个对象注册为可访问的服务，之后该对象的公开方法就能够以远程的方式提供访问。一个 RPC 服务端可以注册多个不同类型的对象，但不允许注册同一类型的多个对象。

一个对象中只有满足如下这些条件的方法，才能被 RPC 服务端设置为可供远程访问：

- 必须是在对象外部可公开调用的方法（首字母大写）；
- 必须有两个参数，且参数的类型都必须是包外部可以访问的类型或者是 Go 内建支持的类型；
- 第二个参数必须是一个指针；
- 方法必须返回一个error类型的值。

以上 4 个条件，可以简单地用如下一行代码表示：

```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

在上面这行代码中，类型 T、T1 和 T2 默认会使用 Go 内置的 encoding/gob 包进行编码解码。关于 encoding/gob 包的内容，稍后我们将会对其进行介绍。

该方法（MethodName）的第一个参数表示由 RPC 客户端传入的参数，第二个参数表示要返回给 RPC 客户端的结果，该方法最后返回一个 error 类型的值。

RPC 服务端可以通过调用 rpc.ServeConn 处理单个连接请求。多数情况下，通过 TCP 或是 HTTP 在某个网络地址上进行监听来创建该服务是个不错的选择。

在 RPC 客户端，Go 的 net/rpc 包提供了便利的 rpc.Dial() 和 rpc.DialHTTP() 方法来与指定的 RPC 服务端建立连接。在建立连接之后，Go 的 net/rpc 包允许我们使用同步或者异步的方式接收 RPC 服务端的处理结果。

调用 RPC 客户端的 Call() 方法则进行同步处理，这时候客户端程序按顺序执行，只有接收完 RPC 服务端的处理结果之后才可以继续执行后面的程序。

当调用 RPC 客户端的 Go() 方法时，则可以进行异步处理，RPC 客户端程序无需等待服务端的结果即可执行后面的程序，而当接收到 RPC 服务端的处理结果时，再对其进行相应的处理。

无论是调用 RPC 客户端的 Call() 或者是 Go() 方法，都必须指定要调用的服务及其方法名称，以及一个客户端传入参数的引用，还有一个用于接收处理结果参数的指针。

如果没有明确指定 RPC 传输过程中使用何种编码解码器，默认将使用 Go 标准库提供的 encoding/gob 包进行数据传输。

接下来，我们来看一组 RPC 服务端和客户端交互的示例程序。下面的代码是 RPC 服务端程序。

具体目录如下：

```text
code/011/rpc_protocol
├── client
│   ├── asynchronous    // 异步调用远程RPC服务
│   │   └── client.go
│   └── synchronous     // 同步调用远程RPC服务
│       └── client.go
├── server              // 远程RPC服务
│   └── server.go
└── service             // 远程RPC服务定义
    └── service.go
```

#### service\/service.go

```go
package service

import (
	"errors"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Ardith int

func (t *Ardith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Ardith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}

	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
```

#### server\/server.go

注册服务对象并开启该 RPC 服务的代码如下：

```go
package main

import (
	"code-snippet/code/011/rpc_protocol/service"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	exit := make(chan string)
	ardith := new(service.Ardith)
	err := rpc.Register(ardith)
	if err != nil {
		log.Fatal(err.Error())
	}

	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err.Error())
	}

	go http.Serve(listener, nil)

	<-exit
}
```

此时，RPC 服务端注册了一个 Arith 类型的对象及其公开方法 Arith.Multiply() 和 Arith.Divide() 供 RPC 客户端调用。

#### synchronous\/client.go

RPC 在调用服务端提供的方法之前，必须先与 RPC 服务端建立连接，RPC 客户端可以调用服务端提供的方法。以下是同步方式进行调用，如下列代码所示：

```go
package main

import (
	"code-snippet/code/011/rpc_protocol/service"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal(err.Error())
	}

	args := &service.Args{A: 7, B: 8}
	var reply int
	err = client.Call("Ardith.Multiply", args, &reply)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("ardith: %d * %d = %d\n", args.A, args.B, reply)
}
```

执行代码后结果如下：

```go
ardith: 7 * 8 = 56
```

#### asynchronous\/client.go

RPC 在调用服务端提供的方法之前，必须先与 RPC 服务端建立连接，RPC 客户端可以调用服务端提供的方法。以下是异步方式进行调用，如下列代码所示：

```go
package main

import (
	"code-snippet/code/011/rpc_protocol/service"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal(err.Error())
	}

	args := &service.Args{A: 17, B: 8}
	quotient := new(service.Quotient)
	call := client.Go("Ardith.Divide", args, &quotient, nil)
	<-call.Done

	fmt.Printf("ardith: %d / %d = %d, %d %% %d = %d\n", args.A, args.B, quotient.Quo, args.A, args.B, quotient.Rem)
}
```

执行代码后结果如下：

```go
ardith: 17 / 8 = 2, 17 % 8 = 1
```
