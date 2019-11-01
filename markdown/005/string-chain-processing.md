### Go语言字符串的链式处理——操作与数据分离的设计技巧

使用 SQL 语言从数据库中获取数据时，可以对原始数据进行排序（sort by）、分组（group by）和去重（distinct）等操作，SQL 将数据的操作与遍历过程作为两个部分进行隔离，这样操作和遍历过程就可以各自独立地进行设计，这就是常见的数据与操作分离的设计。

对数据的操作进行多步骤的处理被称为链式处理，本例中使用多个字符串作为数据集合，然后对每个字符串进行一系列的处理，用户可以通过系统函数或者自定义函数对链式处理中的每个环节进行自定义。

首先给出本节完整代码：

```go
package main

import (
	"fmt"
	"strings"
)

// 字符串处理函数，传入字符串切片和处理链
func StringProcess(list []string, chain []func(string) string) []string {
	for index, str := range list {
		result := str
		for _, function := range chain {
			result = function(result)
		}
		list[index] = result
	}
	return list
}

func main() {
	// 待处理的字符串列表
	list := []string{
		"go scanner",
		"go parser",
		"go compiler",
		"go printer",
		"go formater",
	}

	// 函数处理链
	chains := []func(string) string{
		removeStringPrefix,
		strings.TrimSpace,
		strings.ToUpper,
	}

	// 处理字符串
	list = StringProcess(list, chains)

	// 输出处理好的字符串
	for _, str := range list {
		fmt.Println(str)
	}
}

// 自定义的移除前缀的处理函数
func removeStringPrefix(in string) (out string) {
	return strings.TrimPrefix(in, "go")
}
```

具体的运行结果如下所示：

```text
SCANNER
PARSER
COMPILER
PRINTER
FORMATER
```