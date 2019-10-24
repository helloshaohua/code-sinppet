### Go语言修改字符串

枚举在 C# 中是一个独立的类型，可以通过枚举值获取该值对应的字符串。例如，C# 中 Week 枚举值 Monday 为 1，那么可以通过 Week.Monday.ToString() 函数获得 Monday 字符串。

Go语言中也可以实现这一功能，代码如下所示：

转换字符串：

```go
package main

import "fmt"

// 声明芯片类型
type ChipType int

const (
	None ChipType = iota
	CPU	// 中央处理器
	GPU	// 图片处理器
)

func (c ChipType) String() string {
	switch c {
	case None:
		return "None"
	case CPU:
		return "CPU"
	case GPU:
		return "GPU"
	}
	return "N/A"
}

func main() {
	// 输出CPU的值并以字符串和整型格式显示
	fmt.Printf("%s %d", CPU, CPU)
}
```

运行结果：

```text
CPU 1
```

代码说明如下：
第 6 行，将 int 声明为 ChipType 芯片类型。
第 9 行，将 const 里定义的常量值设为 ChipType 类型，且从 0 开始，每行值加 1。
第 14 行，定义 ChipType 类型的方法 String()，返回值为字符串类型。
第 15～22 行，使用 switch 语句判断当前的 ChitType 类型的值，返回对应的字符串。
第 28 行，按整型的格式输出 CPU 的值。