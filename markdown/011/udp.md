## Go语言UDP网络程序设计

UDP 和上一节《TCP网络程序设计》中的 TCP 一样，也工作在网络传输层，但和 TCP 不同的是，它提供不可靠的通信服务。UDP 网络编程也为 C/S(客户端/服务器) 模式，要设计一个服务器，一个或多个客户端。

另外，UDP 是不保证可靠性的通信协议，所以客户端和服务器之间只要建立连接，就可以直接通信，而不用调用 Aceept() 进行连接确认。本节将详细讲解 UDP 网络编程服务器、客户端的设计原理和设计过程。

### UDPAddr 地址结构体

在进行 UDP 网络编程时，服务器或客户端的地址使用 UDPAddr 地址结构体表示，形式如下：

```go
// UDPAddr represents the address of a UDP end point.
type UDPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone
}
```

函数 ResolveUDPAddr() 可以把网络地址转换为 UDPAddr 地址结构，该函数原型定义如下：

```go
// ResolveUDPAddr returns an address of UDP end point.
//
// The network must be a UDP network name.
//
// If the host in the address parameter is not a literal IP address or
// the port is not a literal port number, ResolveUDPAddr resolves the
// address to an address of UDP end point.
// Otherwise, it parses the address as a pair of literal IP address
// and port number.
// The address parameter can use a host name, but this is not
// recommended, because it will return at most one of the host name's
// IP addresses.
//
// See func Dial for a description of the network and address
// parameters.
func ResolveUDPAddr(network, address string) (*UDPAddr, error)
```

在调用函数 ResolveUDPAddr() 时，参数 net 是网络协议名，可以是“udp”、“udp4”或“udp6”。参数 addr 是 IP 地址或域名，如果是 IPv6 地址则必须使用“[]”括起来。另外，端口号以“:”的形式跟随在 IP 地址或域名的后面，端口是可选的。

函数 ResolveUDPAddr() 调用成功后返回一个指向 UDPAddr 结构体的指针，否则返回一个错误类型。

另外，UDPAddr 地址对象还有两个方法：Network() 和 String()，Network() 方法用于返回 UDPAddr 地址对象的网络协议名，比如“udp”；String() 方法可以将 UDPAddr 地址转换成字符串形式。这两个方法原型定义如下：

```go
func (a *UDPAddr) Network() string
func (a *UDPAddr) String() string
```

### UDPConn 对象

在进行 UDP 网络编程时，客户端和服务器之间是通过 UDPConn 对象实现连接的，UDPConn 是 Conn 接口的实现。UDPConn 对象绑定了服务器的网络协议和地址信息。UDPConn 对象定义如下：

```go
// UDPConn is the implementation of the Conn and PacketConn interfaces
// for UDP network connections.
type UDPConn struct {
	conn
}
```

通过 UDPConn 连接对象在客户端和服务器之间进行通信，UDP 并不能保证通信的可靠性和有序性，这些都要由程序员来处理。为此，TCPConn 对象提供了 ReadFromUDP() 方法和 WriteToUDP() 方法，这两个方法直接使用远端主机地址进行数据发送和接收，即便在链路失效的情况下，通信操作都能正常进行。

ReadFromUDP() 方法和 WriteToUDP() 方法的原型定义如下：

```go
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error)
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)
```

ReadFromUDP() 方法调用成功后返回接收字节数和发送方地址，否则返回一个错误类型；WriteToUDP() 方法调用成功后返回发送字节数，否则返回一个错误类型。

### UDP 服务器设计

在 UDP 网络编程中，服务器工作过程如下：

1.UDP 服务器首先注册一个公知端口，然后调用 ListenUDP() 函数在这个端口上创建一个 UDPConn 连接对象，并在该对象上和客户端建立不可靠连接。
1.如果服务器和某个客户端建立了 UDPConn 连接，就可以使用该对象的 ReadFromUDP() 方法和 WriteToUDP() 方法相互通信了。
1.不管上一次通信是否完成或正常，UDP 服务器依然会接受下一次连接请求。

函数 ListenUDP() 原型定义如下：

```go
// ListenUDP acts like ListenPacket for UDP networks.
//
// The network must be a UDP network name; see func Dial for details.
//
// If the IP field of laddr is nil or an unspecified IP address,
// ListenUDP listens on all available IP addresses of the local system
// except multicast IP addresses.
// If the Port field of laddr is 0, a port number is automatically
// chosen.
func ListenUDP(network string, laddr *UDPAddr) (*UDPConn, error)
```

在调用函数 ListenUDP() 时，参数 network 是网络协议名，可以是“udp”、“udp4”或“udp6”。参数 laddr 是服务器本地地址，可以是任意活动的主机地址，或者是内部测试地址“127.0.0.1”。该函数调用成功，返回一个 UDPConn 对象；调用失败，返回一个错误类型。

#### 【示例 1】UDP Server 端设计，服务器使用本地地址，服务端口号为 8826。

服务器设计工作模式采用循环服务器，对每一个连接请求调用线程 handleClient 来处理。

```go
// UDP Server 端设计
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 获取 UDPAddr 结构体指针
	udpAddr, err := net.ResolveUDPAddr("udp", ":8826")
	checkError(err)

	// 监听服务端口
	udpConn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	// 循环处理每一个客户端
	for {
		handleClient(udpConn)
	}
}

func handleClient(udpConn *net.UDPConn) {
	bytes := make([]byte, 512)
	n, udpAddr, err := udpConn.ReadFromUDP(bytes)
	if err != nil {
		return
	}

	fmt.Printf("client address: %s, received request data: %s\n", udpAddr.String(), string(bytes[0:n]))
	udpConn.WriteToUDP([]byte("welcome client"), udpAddr)
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
```

### UDP 客户端设计

在 UDP 网络编程中，客户端工作过程如下：

1) UDP 客户端在获取了服务器的服务端口号和服务地址之后，可以调用 DialUDP() 函数向服务器发出连接请求，如果请求成功会返回 UDPConn 对象。

2) 客户端可以直接调用 UDPConn 对象的 ReadFromUDP() 方法或 WriteToUDP() 方法，与服务器进行通信活动。

3) 通信完成后，客户端调用 Close() 方法关闭 UDPConn 连接，断开通信链路。

函数 DialUDP() 原型定义如下：

```go
// DialUDP acts like Dial for UDP networks.
//
// The network must be a UDP network name; see func Dial for details.
//
// If laddr is nil, a local address is automatically chosen.
// If the IP field of raddr is nil or an unspecified IP address, the
// local system is assumed.
func DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error) 
```

在调用函数 DialUDP() 时，参数 network 是网络协议名，可以是“udp”、“udp4”或“udp6”。参数 laddr 是本地主机地址，可以设为 nil。参数 raddr 是对方主机地址，必须指定不能省略。函数调用成功后，返回 UDPConn 对象；调用失败，返回一个错误类型。

方法 Close() 的原型定义如下：

```go
func (c *UDPConn) Close() error
```

该方法调用成功后，关闭 UDPConn 连接；调用失败，返回一个错误类型。

#### 【示例 2】UDP Client 端设计，客户端通过内部测试地址“127.0.0.1”和端口 8826 和服务器建立通信连接。


````go
// UDP Client端设计
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port\n", os.Args[0])
		os.Exit(1)
	}

	// 服务IP和端口号
	service := os.Args[1]

	// 获取 UDPAddr 结构体指针
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	// 连接到UDP网络服务器
	udpConn, err := net.DialUDP("udp", nil, udpAddr)

	fmt.Printf("%T\n", udpAddr)

	// 发送数据到UDP网络服务器
	_, err = udpConn.Write([]byte("hello server"))
	checkError(err)

	// 读取UDP网络服务器响应数据
	bytes := make([]byte, 512)
	n, addr, err := udpConn.ReadFromUDP(bytes)
	checkError(err)
	fmt.Printf("server address: %s, recevied response data: %s\n", addr.String(), string(bytes[0:n]))

	udpConn.Close()
	os.Exit(0)
}

// checkError 错误检测，存在错误时退出
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s\n", err.Error())
		os.Exit(1)
	}
}
````

编译并运行服务器端和客户端，测试过程如下：

通过命令终端1启动服务器端程序

```shell
# 启动服务器端
$ go run server.go
```

通过命令终端2启动客户端程序

```shell
# 启动客户端
$ go run client.go localhost:8826
```

终端1服务器端程序接收到请求信息如下：

```shell
# 启动服务器端
$ go run server.go
client address: 127.0.0.1:55350, received request data: hello server
```

终端2客户端程序接收到响应信息如下：

```shell
# 启动客户端
$ go run client.go localhost:8826
server address: 127.0.0.1:8826, recevied response data: welcome client
```

![运行服务器端和客户端](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191124_18.gif)

通过测试结果会发现，采用 TCP 时必须先启动服务器，然后才能正常启动客户端，如果服务器中断，则客户端也会异常退出。而采用 UDP 时，客户端和服务器启动没有先后次序，而且即便是服务器异常退出，客户端也能正常工作。

总之，TCP 可以保证客户端、服务器双方按照可靠有序的方式进行通信，但通信效率低；而 UDP 虽然不能保证通信的可靠性，但通信效率要高得多，在有些场合还是非常有用的。