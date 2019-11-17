### Go语言使用reflect.Type显示一个类型的方法集

本节通过示例来演示如何使用 reflect.Type 来打印任意值的类型和枚举它的方法：

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Print 打印值x的方法集
func Print(x interface{}) {
	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)
	fmt.Printf("type %s\n", t)
	for i := 0; i < v.NumMethod(); i++ {
		mt := v.Method(i).Type()
		fmt.Printf("func (%s) %s %s\n", t, t.Method(i).Name, strings.TrimPrefix(mt.String(), "func"))
	}
	fmt.Println()
}

func main() {
	Print(time.Hour)
	Print(new(strings.Replacer))
}
```

reflect.Type 和 reflect.Value 都提供了一个 Method 方法。每次 t.Method(i) 调用将一个 reflect.Method 的实例，对应一个用于描述一个方法的名称和类型的结构体。每次 v.Method(i) 方法调用都返回一个 reflect.Value 以表示对应的值，也就是一个方法是帮到它的接收者的。

使用 reflect.Value.Call 方法，将可以调用一个 Func 类型的 Value，但是这个例子中只用到了它的类型。这是属于 time.Duration 和 *strings.Replacer 两个类型的方法：

```text
type time.Duration
func (time.Duration) Hours () float64
func (time.Duration) Microseconds () int64
func (time.Duration) Milliseconds () int64
func (time.Duration) Minutes () float64
func (time.Duration) Nanoseconds () int64
func (time.Duration) Round (time.Duration) time.Duration
func (time.Duration) Seconds () float64
func (time.Duration) String () string
func (time.Duration) Truncate (time.Duration) time.Duration

type *strings.Replacer
func (*strings.Replacer) Replace (string) string
func (*strings.Replacer) WriteString (io.Writer, string) (int, error)
```