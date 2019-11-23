### Go语言Cookie的设置与读取

Web 开发中一个很重要的议题就是如何做好用户整个浏览过程的控制，因为 HTTP 协议是无状态的，所以用户的每一次请求都是无状态的，不知道在整个 Web 操作过程中哪些连接与该用户有关。应该如何来解决这个问题呢？Web 里面经典的解决方案是 Cookie 和 Session。

Cookie 机制是一种客户端机制，把用户数据保存在客户端，而 Session 机制是一种服务器端的机制，服务器使用一种类似于散列表的结构来保存信息，每一个网站访客都会被分配给一个唯一的标识符，即 sessionID。

sessionID 的存放形式无非两种：要么经过 URL 传递，要么保存在客户端的 Cookie 里。当然，也可以将 Session 保存到数据库里，这样会更安全，但效率方面会有所下降。本节主要介绍 Go语言使用 Cookie 的方法。

#### 设置 Cookie

Go语言中通过 net/http 包中的 SetCookie 来设置 Cookie：

```go
http.SetCookie(writer ResponseWriter, cookie *Cookie)
```

writer 表示需要写入的 response，cookie 是一个 struct，具体结构体定义如下：

```go
// A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an
// HTTP response or the Cookie header of an HTTP request.
// 
// Cookie表示在HTTP响应的Set-Cookie头或HTTP请求的Cookie头中发送的HTTP Cookie。
// 
// See https://tools.ietf.org/html/rfc6265 for details.
type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified(意味着没有指定 Max-Age 的值).
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0' (意味着现在就删除 Cookie，等价于 Max-Age=0)
	// MaxAge>0 means Max-Age attribute present and given in seconds (意味着 Max-Age 属性存在并以秒为单位存在)
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs(未解析的 attribute-value 属性位对)
}
```

下面来看一个如何设置 Cookie 的例子：

```go
package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// 设置 Cookie
	http.HandleFunc("/set-cookie", func(writer http.ResponseWriter, request *http.Request) {
		http.SetCookie(writer, &http.Cookie{
			Name:     "username",
			Value:    base64.StdEncoding.EncodeToString([]byte("wumoxi")),
			Expires:  time.Now().AddDate(0, 1, 0),
			HttpOnly: true,
		})
		io.WriteString(writer, "hello world")
		return
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

通过浏览器的开发者工作可以查看 Cookie 列表：

![设置 Cookie](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191123_4.png)
 
#### 读取 Cookie

上面的例子演示了如何设置 Cookie 数据，这里演示如何读取 Cookie：

```go
package main

import (
	"encoding/base64"
	"log"
	"net/http"
)

func main() {
	// 获取 Cookie
	http.HandleFunc("/get-cookie", func(writer http.ResponseWriter, request *http.Request) {
		// 获取 Cookie
		cookie, err := request.Cookie("username")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		// 解析Cookie值
		bytes, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// 写入响应对象
		_, err = writer.Write(bytes)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		return
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

通过浏览器查看获取到的 Cookie 值：

![读取 Cookie](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191123_7.png)

#### 读取全部 Cookie

还有另外一种读取方式，获取所有Cookies。

```go
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 获取所有 Cookie
	http.HandleFunc("/get-cookie-all", func(writer http.ResponseWriter, request *http.Request) {
		// 获取所有的 Cookie
		cookies := request.Cookies()
		for _, cookie := range cookies {
			bytes, _ := base64.StdEncoding.DecodeString(cookie.Value)
			fmt.Fprintf(writer, "cookie name: %s, cookie value: %s\n", cookie.Name, string(bytes))
		}
		return
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

通过浏览器查看所有获取到的 Cookie 信息：

![读取 Cookie](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191123_8.png)

可以看到通过 request 获取 Cookie 非常方便。