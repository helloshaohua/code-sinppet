## Go语言TCP网络程序设计

TCP 是工作在TPC/IP网络模型中的传输层协议，它是一种面向连接的可靠的通信协议。TCP 网络程序设计属于 C/S(客户端/服务器) 模式，一般要设计一个服务器程序，一个或多个客户端程序。另外，TCP 是面向连接的通信协议，所以客户端要和服务器进行通信，首先要在通信双方之间建立通信连接。本节将详细讲解 TCP 网络编程中服务器、客户端的设计原理和设计过程。

### TCPAddr 地址结构体

```go
// TCPAddr represents the address of a TCP end point(表示TCP端点的地址).
type TCPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone(IPv6范围的寻址区)
}
```

函数 ResolveTCPAddr() 可以把网络地址转换为 TCPAddr 地址结构，该函数原型定义如下：

```go
func ResolveTCPAddr(network, address string) (*TCPAddr, error)
```

在调用函数 ResolveTCPAddr() 时，参数 network 是网络协议名，可以是“tcp”、“tcp4”或“tcp6”。参数 address 是 IP 地址或域名，如果是 IPv6 地址则必须使用“[]”括起来。另外，端口号以“:”的形式跟随在 IP 地址或域名的后而，端口是可选的。例如：“www.google.com:80”或“127.0.0.1:21”。

还有一种特例，就是对于 HTTP 服务器，当主机地址为本地测试地址时 (127.0.0.1)，可以直接使用端口号作为 TCP 连接地址，形如“:80”。

函数 ResolveTCPAddr() 调用成功后返回一个指向 TCPAddr 结构体的指针，否则返回一个错误类型。

另外，TCPAddr 地址对象还有两个方法：Network() 和 String()，Network() 方法用于返回 TCPAddr 地址对象的网络协议名，比如“tcp”；String() 方法可以将 TCPAddr 地址转换成字符串形式。这两个方法原型定义如下：

```go
// Network returns the address's network name, "tcp"(返回地址的网络名称,如 tcp).
func (a *TCPAddr) Network() string
// String 将 TCPAddr 地址转换成字符串形式
func (a *TCPAddr) String() string
```

####【示例 1】TCP 连接地址

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s network address\n", os.Args[0])
		os.Exit(1)
	}

	// 网络类型
	network := os.Args[1]

	// 域名或IP
	address := os.Args[2]

	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.ResolveTCPAddr error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "network: %s, address: %s\n", tcpAddr.Network(), tcpAddr.String())
	os.Exit(0)
}
```

编译并运行该程序，测试过程如下：

```shell
$ go run tcp-connect-address.go tcp localhost:8080
network: tcp, address: 127.0.0.1:8080
```

### TCPConn 对象

在进行 TCP 网络编程时，客户端和服务器之间是通过 TCPConn 对象实现连接的，TCPConn 是 Conn 接口的实现。TCPConn 对象绑定了服务器的网络协议和地址信息，TCPConn 对象定义如下：

```go
// TCPConn is an implementation of the Conn interface for TCP network connections(TCPConn是用于TCP网络连接的Conn接口的实现).
type TCPConn struct {
	conn
}
```

通过 TCPConn 连接对象，可以实现客户端和服务器间的全双工通信。可以通过 TCPConn 对象的 Read() 方法和 Write() 方法，在服务器和客户端之间发送和接收数据。Read() 方法和 Write() 方法的原型定义如下：

```go
// Read implements the Conn Read method.
func (c *conn) Read(b []byte) (int, error) 

// Write implements the Conn Write method.
func (c *conn) Write(b []byte) (int, error) 
```

因为 `conn` 是 `TCPConn` 内嵌结构体，所以 `TCPConn` 也继承了 `Read` 和 `Write` 方法。Read() 方法调用成功后会返回接收到的字节数，调用失败返回一个错误类型；Write() 方法调用成功后会返回正确发送的字节数，调用失败返回一个错误类型。另外，这两个方法的执行都会引起阻塞。

### TCP 服务器设计

前面讲了 Go语言网络编程和传统 Socket 网络编程有所不同，TCP 服务器的工作过程如下：

1. TCP 服务器首先注册一个公知端口，然后调用 ListenTCP() 函数在这个端口上创建一个 TCPListener 监听对象，并在该对象上监听客户端的连接请求。
2. 启用 TCPListener 对象的 Accept() 方法接收客户端的连接请求，并返回一个协议相关的 Conn 对象，这里就是 TCPConn 对象。
3. 如果返回了一个新的 TCPConn 对象，服务器就可以调用该对象的 Read() 方法接收客户端发来的数据，或者调用 Write() 方法向客户端发送数据了。

ListenTCP() 函数、TCPListener 对象的原型定义如下：

```go
func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)

// TCPListener is a TCP network listener. Clients should typically
// use variables of type Listener instead of assuming TCP.
type TCPListener struct {
	fd *netFD
	lc ListenConfig
}
```

在调用函数 ListenTCP() 时，参数 network 是网络协议名，可以是“tcp”、“tcp4”或“tcp6”。参数 laddr 是服务器本地地址，可以是任意活动的主机地址，或者是内部测试地址“127.0.0.1”。该函数调用成功，返回一个 TCPListener 对象；调用失败，返回一个错误类型。

TCPListener 对象的 AcceptTCP() 方法原型定义如下：

```go
// AcceptTCP accepts the next incoming call and returns the new
// connection.
func (l *TCPListener) AcceptTCP() (*TCPConn, error)
``` 

AcceptTCP() 方法调用成功后，返回 TCPConn 对象；否则，返回一个错误类型。

服务器和客户端的通信连接建立成功后，就可以使用 Read() 和 Write() 方法收发数据。在通信过程中，如果还想获取通信双方的地址信息，可以使用 LocalAddr() 方法和 RemoteAddr() 方法来完成，这两个方法原型定义如下：

```go
// LocalAddr returns the local network address.
// The Addr returned is shared by all invocations of LocalAddr, so
// do not modify it.
func (c *conn) LocalAddr() Addr

// RemoteAddr returns the remote network address.
// The Addr returned is shared by all invocations of RemoteAddr, so
// do not modify it.
func (c *conn) RemoteAddr() Addr
```

LocalAddr() 方法会返回本地主机地址，RemoteAddr() 方法返回远端主机地址。


####【示例 2】TCP Server 端设计，服务器使用本地地址，服务端口号为 8825。

服务器设计工作模式采用循环服务器，对每一个连接请求调用线程 handleClient 来处理。


```go
// TCP Server 端设计
package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// 获取TCPAddr结构体指针对象
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8825")
	checkError(err)

	// 监听服务端口号
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// 循环接收连接请求及请求处理
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 客户端请求处理
		handleClient(tcpConn)
	}
}

// 客户端请求处理
func handleClient(conn *net.TCPConn) {
	defer conn.Close()
	buffer := make([]byte, 512)

	for {
		// 读取请求内容
		i, err := conn.Read(buffer[0:])
		if err != nil {
			return
		}

		// 获取客户端IP
		clientIP := conn.RemoteAddr()
		fmt.Printf("client address: %s, received request data: %s\n", clientIP.String(), string(buffer[0:i]))

		// 写入响应内容
		_, err = conn.Write([]byte("welcome client"))
		if err != nil {
			return
		}
	}
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
```

### TCP 客户端设计

在 TCP 网络编程中，客户端的工作过程如下：

1.TCP 客户端在获取了服务器的服务端口号和服务地址之后，可以调用 DialTCP() 函数向服务器发出连接请求，如果请求成功会返回 TCPConn 对象。
1.客户端调用 TCPConn 对象的 Read() 或 Write() 方法，与服务器进行通信活动。
1.通信完成后，客户端调用 Close() 方法关闭连接，断开通信链路。

DialTCP() 函数原型定义如下：

```go
// DialTCP acts like Dial for TCP networks.
//
// The network must be a TCP network name; see func Dial for details.
//
// If laddr is nil, a local address is automatically chosen.
// If the IP field of raddr is nil or an unspecified IP address, the
// local system is assumed.
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
```

在调用函数 DialTCP() 时，参数 network 是网络协议名，可以是“tcp”、“tcp4”或“tcp6”。参数 laddr 是本地主机地址，可以设为 nil。参数 raddr 是对方主机地址，必须指定不能省略。函数调用成功后，返回 TCPConn 对象；调用失败，返回一个错误类型。

方法 Close() 的原型定义如下：

```go
// Close closes the connection.
func (c *conn) Close() error
```

该方法调用成功后，关闭 TCPConn 连接；调用失败，返回一个错误类型。

#### 【示例 3】TCP Client 端设计，客户端通过内部测试地址“127.0.0.1”和端口 8825 和服务器建立通信连接。

```go
// TCP Client端设计
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 命令行参数检测
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port\n", os.Args[0])
		os.Exit(1)
	}

	// 获取命令行参数域名或IP地址和端口号
	service := os.Args[1]

	// 获取 TCPAddr 结构体指针对象
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	// 连接到TCP
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer tcpConn.Close()

	// 获取服务器IP
	serverIP := tcpConn.RemoteAddr()

	// 发送请求数据
	n, err := tcpConn.Write([]byte("hello server"))
	checkError(err)

	// 读取响应数据
	bytes := make([]byte, 512)
	n, err = tcpConn.Read(bytes[0:])
	checkError(err)
	fmt.Fprintf(os.Stdout, "server address: %s, recevied response data: %s\n", serverIP.String(), string(bytes[0:n]))

	// 正常退出程序
	os.Exit(0)
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
```

编译并运行服务器端和客户端，测试过程如下：

通过命令终端1启动服务器端程序

```shell
# 启动服务器端
$ go run server.go
```

通过命令终端2启动客户端程序

```shell
# 启动客户端
$ go run client.go localhost:8825
```

终端1服务器端程序接收到请求信息如下：

```shell
# 启动服务器端
$ go run server.go
client address: 127.0.0.1:64809, received request data: hello server
```

终端2客户端程序接收到响应信息如下：

```shell
# 启动客户端
$ go run client.go localhost:8825
server address: 127.0.0.1:8825, recevied response data: welcome client
```

![运行服务器端和客户端](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191124_17.gif)

从上述测试结果可以看出，服务器注册了一个公知端口 8825，而当客户端与服务器建立连接后，客户端会生成一个临时端口“50813”与服务器进行通信。服务器不管启动多少次端口号都是 8825，而客户端每一次重新启动端口号都不一样。


### 使用 Goroutine 实现并发服务器

前面的讲解中服务器设计采用循环服务器设计模式，这种服务器设计简单但缺陷明显。因为这种服务器一旦启动，就一直阻塞监听客户端的连接请求，直至服务器关闭。所以，循环服务器很耗费系统资源。

解决问题的方法是采用并发服务器模式，在这种模式中，对每一个客户端的连接请求，服务器都会创建一个新的进程、线程或者协程进行响应，而服务器还可以去处理其他任务。Goroutine 即协程是一种比线程更轻量级的任务单位，所以这里就使用 Goroutine 来实现并发服务器的设计。

#### 【示例 4】并发服务器设计，服务器使用本地地址，服务端口号为 8825。

服务器设计工作模式采用并发服务器模式，对每一个连接请求创建一个能调用 handleClient() 函数的 Goroutine 来处理。

```go
// TCP Server 端设计
package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// 获取TCPAddr结构体指针对象
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8825")
	checkError(err)

	// 监听服务端口号
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// 循环接收连接请求及请求处理
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 使用goroutine并发处理客户端请求
		go handleClient(tcpConn)
	}
}

// 客户端请求处理
func handleClient(conn *net.TCPConn) {
	// 逆序调用 Close() 保证连接能正常关闭
	defer conn.Close()
	buffer := make([]byte, 512)

	for {
		// 读取请求内容
		i, err := conn.Read(buffer[0:])
		if err != nil {
			return
		}

		// 获取客户端IP
		clientIP := conn.RemoteAddr()
		fmt.Printf("client address: %s, received request data: %s\n", clientIP.String(), string(buffer[0:i]))

		// 写入响应内容
		_, err = conn.Write([]byte("welcome client"))
		if err != nil {
			return
		}
	}
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
```

编译并运行服务器端和客户端，测试过程如下：


通过命令终端1启动服务器端程序

```shell
# 启动服务器端
$ go run server.go
```

通过命令终端2启动客户端程序

```shell
# 启动客户端
$ go run client.go localhost:8825
```

终端1服务器端程序接收到请求信息如下：

```shell
# 启动服务器端
$ go run server.go
client address: 127.0.0.1:64809, received request data: hello server
```

终端2客户端程序接收到响应信息如下：

```shell
# 启动客户端
$ go run client.go localhost:8825
server address: 127.0.0.1:8825, recevied response data: welcome client
```

![运行服务器端和客户端](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191124_17.gif)

通过测试可以发现，并发服务器可以同时响应多个客户端的连接请求，并能和多个客户端并发通信，尤其在多核心系统平台上，这种通信模式效率更高。而循环服务器只能按客户端的请求队列次序，一个一个地为客户端提供通信服务，通信效率低下。