### 使用指针变量获取命令行的输入信息

Go语言内置的 flag 包实现了对命令行参数的解析，flag 包使得开发命令行工具更为简单。

下面的代码通过提前定义一些命令行指令和对应的变量，并在运行时输入对应的参数，经过 flag 包的解析后即可获取命令行的数据。

【示例】获取命令行输入：

```go
package main

import (
	"flag"
	"fmt"
)

// 定义命令行参数
var mode = flag.String("mode", "", "process mode")

func main() {
	// 解析命令行参数
	flag.Parse()

	// 输出命令行参数
	fmt.Println(*mode)
}
```

将这段代码命名为 main.go，然后使用如下命令行运行：

```shell
go run main.go -mode=fast
```

命令行输出结果如下：

```text
fast
```

代码说明如下：
- 第 9 行，通过 flag.String，定义一个 mode 变量，这个变量的类型是 *string。后面 3 个参数分别如下：
    - 参数名称：在命令行输入参数时，使用这个名称。
    - 参数值的默认值：与 flag 所使用的函数创建变量类型对应，String 对应字符串、Int 对应整型、Bool 对应布尔型等。
    - 参数说明：使用 -help 时，会出现在说明中。
- 第 13 行，解析命令行参数，并将结果写入到变量 mode 中。
- 第 16 行，打印 mode 指针所指向的变量。

由于之前已经使用 flag.String 注册了一个名为 mode 的命令行参数，flag 底层知道怎么解析命令行，并且将值赋给 mode*string 指针，在 Parse 调用完毕后，无须从 flag 获取值，而是通过自己注册的这个 mode 指针获取到最终的值。代码运行流程如下图所示。

![命令行参数与变量的关系](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191023_14.jpg)


