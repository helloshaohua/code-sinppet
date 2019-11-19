### Go语言通过反射访问结构体成员的值

反射值对象（reflect.Value）提供对结构体访问的方法，通过这些方法可以完成对结构体任意值的访问，如下表所示。

| 方  法                                         | 备  注                                                       |
| ---------------------------------------------- | ------------------------------------------------------------ |
| Field(i int) Value                             | 根据索引，返回索引对应的结构体成员字段的反射值对象。当值不是结构体或索引超界时发生宕机 |
| NumField() int                                 | 返回结构体成员字段数量。当值不是结构体或索引超界时发生宕机   |
| FieldByName(name string) Value                 | 根据给定字符串返回字符串对应的结构体字段。没有找到时返回零值，当值不是结构体或索引超界时发生宕机 |
| FieldByIndex(index []int) Value                | 多层成员访问时，根据 []int 提供的每个结构体的字段索引，返回字段的值。 没有找到时返回零值，当值不是结构体或索引超界时发生宕机 |
| FieldByNameFunc(match func(string) bool) Value | 根据匹配函数匹配需要的字段。找到时返回零值，当值不是结构体或索引超界时发生宕机 |

下面代码构造一个结构体包含不同类型的成员。通过 reflect.Value 提供的成员访问函数，可以获得结构体值的各种数据。

反射访问结构体成员值：

```go
package main

import (
	"fmt"
	"reflect"
)

// 定义结构体
type dummy struct {
	a       int
	b       string
	float32        // 嵌入字段
	bool           // 嵌入字段
	next    *dummy // 类型为自身结构体类型
}

func main() {
	// 值包装结构体
	d := reflect.ValueOf(dummy{
		next: &dummy{a: 17995010},
	})

	// 获取字段数量
	fmt.Printf("NumField: %d\n", d.NumField())

	// 获取索引为2的字段(float32字段)
	float32Field := d.Field(2)

	// 输出字段类型
	fmt.Printf("Field Type: %s\n", float32Field.Type())

	// 根据字段名查找字段
	fmt.Printf("FieldByName(\"b\").Type: %s\n", d.FieldByName("b").Type())

	// 根据索引查找值中 next 字段的 int 字段的类型
	fmt.Printf("FieldByIndex([]int{4, 0}).Type(): %s\n", d.FieldByIndex([]int{4, 0}).Type())

	// 根据索引查找值中 next 字段的 int 字段的值
	fmt.Printf("FieldByIndex([]int{4, 0}).Type(): %d\n", d.FieldByIndex([]int{4, 0}).Int())
}
```

代码说明如下：

- 第 9 行，定义结构体，结构体的每个字段的类型都不一样。
- 第 19 行，实例化结构体并包装为 reflect.Value 类型，成员中包含一个 *dummy 的实例。
- 第 24 行，获取结构体的字段数量。
- 第 27 和 30 行，获取索引为2的字段值（float32 字段），并且打印类型。
- 第 33 行，根据b字符串，查找到 b 字段的类型。
- 第 36 行，[]int{4,0} 中的 4 表示，在 dummy 结构中索引值为 4 的成员，也就是 next。next 的类型为 dummy，也是一个结构体，因此使用 []int{4,0} 中的 0 继续在 next 值的基础上索引，结构为 dummy 中索引值为 0 的 a 字段，类型为 int。
- 第 39 行，[]int{4,0} 中的 4 表示，在 dummy 结构中索引值为 4 的成员，也就是 next。next 的类型为 dummy，也是一个结构体，因此使用 []int{4,0} 中的 0 继续在 next 值的基础上索引，结构为 dummy 中索引值为 0 的 a 字段值，类型为 17995010。

代码输出如下：

```text
NumField: 5
Field Type: float32
FieldByName("b").Type: string
FieldByIndex([]int{4, 0}).Type(): int
FieldByIndex([]int{4, 0}).Type(): 17995010
```