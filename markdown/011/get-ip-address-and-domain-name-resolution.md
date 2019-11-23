## Go语言获取IP地址和域名解析

主机地址是网络通信最重要的数据之一，net 包中定义了三种类型的主机地址数据类型：IP、IPMask 和 IPAddr，它们分别用来存储协议相关的网络地址。

### IP 地址类型

在 net 包中，IP 地址类型被定义成一个 byte 型数数组，格式如下：

```go
single-address
// Functions in this package accept either 4-byte (IPv4)
// or 16-byte (IPv6) slices as input.
//
// Note that in this documentation, referring to an
// IP address as an IPv4 address or an IPv6 address
// is a semantic property of the address, not just the
// length of the byte slice: a 16-byte slice can still
// be an IPv4 address.
type IP []byte
```

在 net 包中，有几个函数可以将 IP 地址类型作为函数的返回值类型，比如 ParseIP() 函数，该函数原型定义如下：

```go
func ParseIP(s string) IP
```

ParseIP() 函数的主要作用是分析 IP 地址的合法性，如果是一个合法的 IP 地址，ParseIP() 函数将返回一个 IP 地址对象。如果是一个非法 IP 地址，ParseIP() 函数将返回 nil。

还可以使用 IP 对象的 String() 方法将 IP 地址转换成字符串格式，String() 方法的原型定义如下：

```go
func (ip IP) String() string
```

如果是 IPv4 地址，String() 方法将返回一个点分割十进制格式的 IP 地址，如“192.168.0.1”。如果是 IPv6 地址，String() 方法将返回使用“:”分隔的地址形式，如“2000:0:0:0:0:0:0:1”。另外注意一个特例，对于地址“0:0:0:0:0:0:0:1”的返回结果是省略格式“::1”。

IP 地址类型。

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s (IP address)\n", os.Args[0])
		os.Exit(1)
	}

	// Get the input ip address
	address := os.Args[1]

	// Parse IP
	ip := net.ParseIP(address)
	if ip != nil {
		fmt.Printf("The address is %s.\n", ip.String())
	} else {
		fmt.Println("invalid address")
	}
	os.Exit(0)
}
```

编译并运行该程序，测试过程如下：

```text
从键盘输入：192.168.0.1
输出结果为：The address is 192.168.0.1
从键盘输入：192.168.0.256
输出结果为：Inval id address
从键盘输入：0:0:0:0:0:0:0:1
输出结果为：::1
```

### IPMask 地址类型

在 Go语言中，为了方便子网掩码操作与计算，net 包中还提供了 IPMask 地址类型。在前面讲过，子网掩码地址其实就是一个特殊的 IP 地址，所以 IPMask 类型也是一个 byte 型数组，格式如下：

```go
// An IP mask is an IP address.
type IPMask []byte
```

函数 IPv4Mask() 可以通过一个 32 位 IPv4 地址生成子网掩码地址，调用成功后返回一个 4 字节的十六进制子网掩码地址。IPv4Mask() 函数原型定义如下：

```go
func IPv4Mask(a, b, c, d byte) IPMask
```

另外，还可以使用主机地址对象的 DefaultMask() 方法获取主机默认子网掩码地址，DefaultMask() 方法原型定义如下：

```go
func (ip IP) DefaultMask() IPMask
```

要注意的是，只有 IPv4 地址才有默认子网掩码。如果不是 IPv4 地址，DefaultMask() 方法将返回 nil。不管是通过调用 IPv4Mask() 函数，还是执行 DefaultMask() 方法，获取的子网掩码地址都是十六进制格式的。例如，子网掩码地址“255.255.255.0”的十六进制格式是“ffffff00”。

主机地址对象还有一个 Mask() 方法，执行 Mask() 方法后，会返回 IP 地址与子网掩码地址相“与”的结果，这个结果即是主机所处的网络的“网络地址”。Mask() 方法原型定义如下：

```go
func (ip IP) Mask(mask IPMask) IP
```
还可以通过子网掩码对象的 Size() 方法获取掩码位数 (ones) 和掩码总长度 (bits)，如果是一个非标准的子网掩码地址，则 Size() 方法将返回“0,0”。Size() 方法的原型定义如下：

```go
func (m IPMask) Size() (ones, bits int)
```

子网掩码地址：

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Usage: %s (IP address)\n", os.Args[0])
		os.Exit(1)
	}

	// 获取输入的ip地址
	address := os.Args[1]

	// 解析IP
	ip := net.ParseIP(address)
	if ip == nil {
		fmt.Fprintf(os.Stderr, "invalid address")
		os.Exit(1)
	}

	// 获取 IP 地址默认的子网掩码
	defaultMask := ip.DefaultMask()
	fmt.Printf("Subnet mask is: %s\n", defaultMask)

	// 获取主机所在的网络地址
	network := ip.Mask(defaultMask)
	fmt.Printf("Network address is: %s\n", network)

	// 获取掩码位数和掩码总长度
	ones, bits := defaultMask.Size()
	fmt.Printf("Mask bits: %d, Total bits: %d\n", ones, bits)
}
```

编译并运行该程序，结果如下所示：

```shell
$ go run mask-address.go 192.168.0.1
Subnet mask is: ffffff00
Network address is: 192.168.0.0
Mask bits: 24, Total bits: 32
```

### 域名解析(只获取一个IP地址)

在 net 包中，许多函数或方法调用后返回的是一个指向 IPAddr 结构体的指针，结构体 IPAddr 内只定义了两个字段，格式如下：

```go
// IPAddr represents the address of an IP end point.
type IPAddr struct {
	IP   IP
	Zone string // IPv6 scoped addressing zone
}
```

IPAddr 结构体的主要作用是用于域名解析服务 (DNS)，例如，函数 ResolveIPAddr() 可以通过域名解析主机网络地址。ResolveIPAddr() 函数原型定义如下：

```go
func ResolveIPAddr(network, address string) (*IPAddr, error)
```

在调用 ResolveIPAddr() 函数时，参数 network 表示网络类型，可以是“ip”、“ip4”或“ip6”，参数 address 可以是 IP 地址或域名，如果是 IPv6 地址则必须使用“[]”括起来。ResolveIPAddr() 函数调用成功后返回指向 IPAddr 结构体的指针，调用失败返回错误类型 error。

DNS 域名解析：

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s hostname\n", os.Args[0])
		os.Exit(1)
	}

	// 获取命令输入域名
	hostname := os.Args[1]

	// 通过域名获取 IP
	ipAddr, err := net.ResolveIPAddr("ip", hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Resolvtion error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Resolved address is: %s\n", ipAddr.String())
	os.Exit(0)
}
```

编译并运行该程序，结果如下所示：

```shell
$ go run domain-parse.go
Resolved address is: 39.156.69.79
```

还有一些域名配置多个IP地址，那么可以通过 http.LookupHost() 函数。

### 域名解析(获取多个IP地址)

通过 http.LookupHost() 函数解析域名，获取多个绑定到该域名的IP 地址。以下是示例：

```go
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s hostname\n", os.Args[0])
		os.Exit(1)
	}

	// 获取命令输入域名
	hostname := os.Args[1]

	// 通过域名获取 IP
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LookupHost error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Lookup host address is: %s\n", strings.Join(addrs, ", "))
	os.Exit(0)
}

```

编译并运行该程序，结果如下所示：

```shell
$ go run multi-address.go baidu.com
Lookup host address is: 220.181.38.148, 39.156.69.79
```