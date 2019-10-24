### Go语言遍历字符串——获取每一个字符串元素

遍历字符串有下面两种写法。

#### 遍历每一个ASCII字符

遍历 ASCII 字符使用 for 的数值循环进行遍历，直接取每个字符串的下标获取 ASCII 字符，如下面的例子所示。

```go
package main

import "fmt"

func main() {
	theme := "狙击 start"
	for i := 0; i < len(theme); i++ {
		fmt.Printf("ascii: %c %d\n", theme[i], theme[i])
	}
}
```

程序输出如下：

```perl
ascii: ç 231
ascii:  139
ascii:  153
ascii: å 229
ascii:  135
ascii: » 187
ascii:   32
ascii: s 115
ascii: t 116
ascii: a 97
ascii: r 114
ascii: t 116
```

这种模式下取到的汉字“惨不忍睹”。由于没有使用 Unicode，汉字被显示为乱码。

#### 按Unicode字符遍历字符串

同样的内容：

```go
package main

import "fmt"

func main() {
	theme := "狙击 start"
	for _, s := range theme {
		fmt.Printf("unicode: %c %d\n", s, s)
	}
}
```

程序输出如下:

```perl
unicode: 狙 29401
unicode: 击 20987
unicode:   32
unicode: s 115
unicode: t 116
unicode: a 97
unicode: r 114
unicode: t 116
```

可以看到，这次汉字可以正常输出了。

#### 总结 

- ASCII 字符串遍历直接使用下标。
- Unicode 字符串遍历用 for range。