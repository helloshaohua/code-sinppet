### Go语言Dial()函数：建立网络连接

Go语言中 Dial() 函数的原型如下：

```go
func Dial(network, address string) (Conn, error)
```

其中 network 参数是网络协议的名字，address 参数是 IP 地址或域名，而端口号以“:”的形式跟随在地址或域名的后面，端口号可选。如果连接成功，返回连接对象，否则返回 error。


我们来看一下几种常见协议的调用方式。

#### 1、TCP 链接

```go
conn, err := net.Dial("tcp", "192.168.0.10:8825")
```

#### 2、UDP 连接

```go
conn, err := net.Dial("udp", "192.168.0.10:8825")
```

#### 3、ICMP 链接（使用协议名称）

```go
conn, err := net.Dial("ip4:icmp", "www.baidu.com")
```

#### 4、ICMP 链接（使用协议编号）

```go
conn, err := net.Dial("ip4:1", "192.168.0.10")
```

这里我们可以通过以下链接查看协议编号的含义：[http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xml](http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xml)。

目前，Dial() 函数支持如下几种网络协议："tcp"、"tcp4"（仅限 IPv4）、"tcp6"（仅限 IPv6）、"udp"、"udp4"（仅限 IPv4）、"udp6"（仅限 IPv6）、"ip"、"ip4"（仅限 IPv4）和"ip6"（仅限 IPv6）。

在成功建立连接后，我们就可以进行数据的发送和接收。发送数据时，使用 conn 的 Write() 成员方法，接收数据时使用 Read() 方法。