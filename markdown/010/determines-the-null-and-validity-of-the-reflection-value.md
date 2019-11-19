### Go语言IsNil()和IsValid()——判断反射值的空和有效性

反射值对象（reflect.Value）提供一系列方法进行零值和空判定，如下表所示。

| 方 法                                          | 说 明                                                        |
| ---------------------------------------------- | ------------------------------------------------------------ |
| IsNil() bool                                   | 返回值是否为 nil。如果值类型不是通道（channel）、函数、接口、map、指针或 切片时发生 panic，类似于语言层的`v== nil`操作 |
| IsValid() bool                                 | 判断值是否有效。 当值本身非法时，返回 false，例如 reflect Value不包含任何值，值为 nil 等。 |

下面的例子将会对各种方式的空指针进行 IsNil() 和 IsValid() 的返回值判定检测。同时对结构体成员及方法查找 map 键值对的返回值进行 IsValid() 判定，参考下面的代码。

反射值对象的零值和有效性判断：

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// *int的空指针
	var a *int
	fmt.Println("var a *int:", reflect.ValueOf(a).IsNil())

	// nil值
	fmt.Println("nil:", reflect.ValueOf(nil).IsValid())

	// *int类型的空指针
	fmt.Println("(*int)(nil):", reflect.ValueOf((*int)(nil)).Elem().IsValid())

	// 实例化一个结构体
	s := struct{}{}

	// 尝试从结构体中查找一个不存在的字段
	fmt.Println("不存在的结构体成员:", reflect.ValueOf(s).FieldByName("").IsValid())

	// 尝试从结构体中查找一个不存在的方法
	fmt.Println("不存在的结构体方法:", reflect.ValueOf(s).MethodByName("").IsValid())

	// 实例化一个 map
	m := make(map[int]int)

	// 尝试从map中查找一个不存在的键
	fmt.Println("不存在的键:", reflect.ValueOf(m).MapIndex(reflect.ValueOf(3)).IsValid())
}
```

代码说明如下：

- 第 10 行，声明一个 *int 类型的指针，初始值为 nil。
- 第 11 行，将变量 a 包装为 reflect.Value 并且判断是否为空，此时变量 a 为空指针，因此返回 true。
- 第 14 行，对 nil 进行 IsValid() 判定（有效性判定），返回 false。
- 第 17 行，(*int)(nil) 的含义是将 nil 转换为 *int，也就是*int 类型的空指针。此行将 nil 转换为 *int 类型，并取指针指向元素。由于 nil 不指向任何元素，*int 类型的 nil 也不能指向任何元素，值不是有效的。因此这个反射值使用 Isvalid() 判断时返回 false。
- 第 20 行，实例化一个结构体。
- 第 23 行，通过 FieldByName 查找 s 结构体中一个空字符串的成员，如成员不存在，IsValid() 返回 false。
- 第 26 行，通过 MethodByName 查找 s 结构体中一个空字符串的方法，如方法不存在，IsValid() 返回 false。
- 第 29 行，实例化一个 map，这种写法与 make 方式创建的 map 等效。
- 第 32 行，MapIndex() 方法能根据给定的 reflect.Value 类型的值查找 map，并且返回查找到的结果。

代码输出如下：

```text
var a *int: true
nil: false
(*int)(nil): false
不存在的结构体成员: false
不存在的结构体方法: false
不存在的键: false
```

IsNil() 常被用于判断指针是否为空；IsValid() 常被用于判定返回值是否有效。
