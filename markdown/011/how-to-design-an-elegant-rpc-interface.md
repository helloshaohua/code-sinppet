### 如何设计优雅的RPC接口

Go语言的 net/rpc 很灵活，它在数据传输前后实现了编码解码器的接口定义。这意味着，开发者可以自定义数据的传输方式以及 RPC 服务端和客户端之间的交互行为。

RPC 提供的编码解码器接口如下：

```go
// A ClientCodec implements writing of RPC requests and
// reading of RPC responses for the client side of an RPC session.
// The client calls WriteRequest to write a request to the connection
// and calls ReadResponseHeader and ReadResponseBody in pairs
// to read responses. The client calls Close when finished with the
// connection. ReadResponseBody may be called with a nil
// argument to force the body of the response to be read and then
// discarded.
// See NewClient's comment for information about concurrent access.
type ClientCodec interface {
	WriteRequest(*Request, interface{}) error
	ReadResponseHeader(*Response) error
	ReadResponseBody(interface{}) error

	Close() error
}

// A ServerCodec implements reading of RPC requests and writing of
// RPC responses for the server side of an RPC session.
// The server calls ReadRequestHeader and ReadRequestBody in pairs
// to read requests from the connection, and it calls WriteResponse to
// write a response back. The server calls Close when finished with the
// connection. ReadRequestBody may be called with a nil
// argument to force the body of the request to be read and discarded.
// See NewClient's comment for information about concurrent access.
type ServerCodec interface {
	ReadRequestHeader(*Request) error
	ReadRequestBody(interface{}) error
	WriteResponse(*Response, interface{}) error

	// Close can be called multiple times and must be idempotent.
	Close() error
}
```

接口 ClientCodec 定义了 RPC 客户端如何在一个 RPC 会话中发送请求和读取响应。客户端程序通过 WriteRequest() 方法将一个请求写入到 RPC 连接中，并通过 ReadResponseHeader() 和 ReadResponseBody() 读取服务端的响应信息。当整个过程执行完毕后，再通过 Close() 方法来关闭该连接。

接口 ServerCodec 定义了 RPC 服务端如何在一个 RPC 会话中接收请求并发送响应。服务端程序通过 ReadRequestHeader() 和 ReadRequestBody() 方法从一个 RPC 连接中读取请求信息。然后再通过 WriteResponse() 方法向该连接中的 RPC 客户端发送响应。当完成该过程后，通过 Close() 方法来关闭连接。

通过实现上述接口，我们可以自定义数据传输前后的编码解码方式，而不仅仅局限于 Gob。

同样，可以自定义 RPC 服务端和客户端的交互行为。实际上，Go 标准库提供的 net/rpc/json 包，就是一套实现了 rpc.ClientCodec 和 rpc.ServerCodec 接口的 JSON-RPC 模块。