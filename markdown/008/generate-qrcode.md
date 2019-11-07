### Go语言生成二维码

二维码作为一种快速的输入手段越来越流行，支付，添加好友，买东西，扫个二维码就可以，非常方便。那么二维码是如何制作生成的呢？我们如何制作自己的二维码呢？

#### 什么是二维码？

二维条码是指在一维条码的基础上扩展出另一维具有可读性的条码，使用黑白矩形图案表示二进制数据，被设备扫描后可获取其中所包含的信息。一维条码的宽度记载着数据，而其长度没有记载数据。二维条码的长度、宽度均记载着数据。

二维条码有一维条码没有的“定位点”和“容错机制”。容错机制在即使没有辨识到全部的条码、或是说条码有污损时，也可以正确地还原条码上的信息。


#### 使用Go语言生成二维码图片

使用Go语言编程时，生成任意内容的二维码是非常方便的，因为我们有 go-qrcode 这个库。该库的源代码托管在 github 上，大家可以从 github 上（[https://github.com/skip2/go-qrcode](https://github.com/skip2/go-qrcode)）下载并使用这个库。

go-qrcode 的使用很简单，假如我们要为百度官网 https://baidu.com 生成一张 256*256 的图片，可以使用如下代码：

```go
package main

import (
	"log"

	"github.com/skip2/go-qrcode"
)

func main() {
	err := qrcode.WriteFile("https://baidu.com", qrcode.Medium, 256, "./baidu.com.png")
	if err != nil {
		log.Fatal(err)
	}
}
```

这样我们运行代码的时候，就在当前目录下，生成一张 256*256 的二维码，扫描后就可以自动跳转到百度。

```go
func WriteFile(content string, level RecoveryLevel, size int, filename string) error
```

WriteFile 函数的原型定义如上，它有几个参数，大概意思如下：

- content 表示要生成二维码的内容，可以是任意字符串。
- level 表示二维码的容错级别，取值有 Low、Medium、High、Highest。
- size 表示生成图片的 width 和 height，像素单位。
- filename 表示生成的文件名路径。

RecoveryLevel(类型其实就是int)，它的定义和常量如下：

```go
// Error detection/recovery capacity.
//
// There are several levels of error detection/recovery capacity. Higher levels
// of error recovery are able to correct more errors, with the trade-off of
// increased symbol size.
type RecoveryLevel int

const (
	// Level L: 7% error recovery.
	Low RecoveryLevel = iota

	// Level M: 15% error recovery. Good default choice.
	Medium

	// Level Q: 25% error recovery.
	High

	// Level H: 30% error recovery.
	Highest
)
```

RecoveryLevel 越高，二维码的容错能力越好。

#### 生成二维码图片字节

有时候我们不想直接生成一个 PNG 文件存储，我们想对 PNG 图片做一些处理，比如缩放了，旋转了，或者网络传输了等，基于此，我们可以使用 Encode 函数，生成一个 PNG 图片的字节流，这样我们就可以进行各种处理了。

```go
func Encode(content string, level RecoveryLevel, size int) ([]byte, error)
```

用法和 WriteFile 函数差不多，只不过返回的是一个 []byte 字节数组，这样我们就可以对这个字节数组进行处理了。

```go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func main() {
	engine := gin.New()
	engine.GET("/code", func(context *gin.Context) {
		bytes, err := qrcode.Encode("https://jd.com", qrcode.Highest, 400)
		if err != nil {
			context.Writer.Write([]byte("Generate QRCODE Failed"))
		}
		context.Writer.Write(bytes)
	})
	log.Fatal(engine.Run(":8295"))
}
```

直接通过浏览器访问 `localhost:8295/code` 可以获取到生成的二维码如下：

![生成的二维码](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191107_3.png)

#### 自定义二维码

除了以上两种快捷方式，go-qrcode 库还为我们提供了对二维码的自定义方式，比如我们可以自定义二维码的前景色和背景色等。qrcode.New 函数可以返回一个 *QRCode，我们可以对 *QRCode 设置，实现对二维码的自定义。

比如我们设置背景色为绿色，前景色为白色的二维码。

```go
package main

import (
	"image/color"
	"log"

	"github.com/skip2/go-qrcode"
)

func main() {
	code, err := qrcode.New("https://imooc.com", qrcode.Highest)
	if err != nil {
		log.Fatal(err)
	}
	code.BackgroundColor = color.RGBA{50, 205, 50, 255}
	code.ForegroundColor = color.White
	code.WriteFile(400, "./imooc.png")
}
```

指定 *QRCode 的 BackgroundColor 和 ForegroundColor 即可。然后调用 WriteFile 方法生成这个二维码文件。

![自定义二维码](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191107_4.png)






