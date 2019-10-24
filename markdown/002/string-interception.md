### Go语言字符串截取（获取字符串的某一段字符）

获取字符串的某一段字符是开发中常见的操作，我们一般将字符串中的某一段字符称做**子串（substring）**。

下面例子中使用 strings.Index() 函数在字符串中搜索另外一个子串，代码如下：

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// 字符串定义
	tracer := "死神来了，死神bye bye"

	// 搜索中文的逗号位置
	comma := strings.Index(tracer, "，")

	// 获得第二个“死神”的位置
	pos := strings.Index(tracer[comma:], "死神")

	// 输出子字符串 "死神bye bye"
	fmt.Println(comma, pos, tracer[comma+pos:])
}
```

#### 代码说明如下

1)、第13行尝试在 tracer 的字符串中搜索中文的逗号，返回的位置存在 comma 变量中，类型是 int，表示从 tracer 字符串开始的 ASCII 码位置。
   
strings.Index() 函数并没有像其他语言一样，提供一个从某偏移开始搜索的功能。不过我们可以对字符串进行切片操作来实现这个逻辑。


2)、第16行中，tracer[comma:] 从 tracer 的 comma 位置开始到 tracer 字符串的结尾构造一个子字符串，返回给 string.Index() 进行再索引。得到的 pos 是相对于 tracer[comma:] 的结果。

comma 逗号的位置是 12，而 pos 是相对位置，值为 3。我们为了获得第二个“死神”的位置，也就是逗号后面的字符串，就必须让 comma 加上 pos 的相对偏移，计算出 15 的偏移，然后再通过切片 tracer[comma+pos:] 计算出最终的子串，获得最终的结果：“死神bye bye”。

#### 总结

字符串索引比较常用的有如下几种方法：

- strings.Index：正向搜索子字符串。
- strings.LastIndex：反向搜索子字符串。
- 搜索的起始位置可以通过切片偏移制作。


