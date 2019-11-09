### Go语言（Pingo）插件化开发

Pingo 是一个用来为Go语言程序编写插件的简单独立库，因为 Go 本身是静态链接的，因此所有插件都以外部进程方式存在。Pingo 旨在简化标准 RPC 包，支持 TCP 和 Unix 套接字作为通讯协议。当前还不支持远程插件，如果有需要，远程插件很快会提供。

#### 创建一个新插件

使用 Pingo 创建一个插件非常简单，首先新建目录，如 "plugins/hello-world" ，然后在该目录下编写 main.go：

```go
package main

import "github.com/dullgiulio/pingo"

// 创建要导出的对象
type MyPlugin struct {
}

// 导出的方法，带有RPC 签名
func (p *MyPlugin) SayHello(name string, msg *string) error {
	*msg = "hello, " + name
	return nil
}

func main() {
	plugin := &MyPlugin{}
	pingo.Register(plugin)
	pingo.Run()
}
```

编译：

```perl
$ cd plugins/hello-world
$ go build
```

#### 调用插件

接下来就可以调用该插件：


```go
package main

import (
	"log"

	"github.com/dullgiulio/pingo"
)

func main() {
	// 从创建的可执行文件中创建一个新插件，通过TCP连接到它
	p := pingo.NewPlugin("tcp", "code/008/plugins/hello-world/hello-world")
	// 启动插件
	p.Start()
	// 使用完插件后停止它
	defer p.Stop()

	resp := ""
	// 从先前创建的对象调用方法
	if err := p.Call("MyPlugin.SayHello", "Go developer", &resp); err != nil {
		log.Print(err)
	} else {
		log.Print(resp)
	}
}
```

代码执行如下：

```perl
2019/11/09 12:38:52 hello, Go developer
```

#### 项目目录

[hello-world](https://github.com/wumoxi/code-sinppet/tree/master/code/008/plugins)

