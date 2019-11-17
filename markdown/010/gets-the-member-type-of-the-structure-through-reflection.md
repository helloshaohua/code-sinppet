### Go语言通过反射获取结构体的成员类型

任意值通过 reflect.TypeOf() 获得反射对象信息后，如果它的类型是结构体，可以通过反射值对象（reflect.Type）的 NumField() 和 Field() 方法获得结构体成员的详细信息。与成员获取相关的 reflect.Type 的方法如下表所示。

| 方法                                         | 说明                                                         |
| :------------------------------------------- | ------------------------------------------------------------ |
| Field(i int) StructField                     | 根据索引，返回索引对应的结构体字段的信息。当值不是结构体或索引超界时发生宕机 |
| NumField() int                               | 返回结构体成员字段数量。当类型不是结构体或索引超界时发生宕机 |
| FieldByName(name string) (StructField, bool) | 根据给定字符串返回字符串对应的结构体字段的信息。没有找到时 bool 返回 false，当类型不是结构体或索引超界时发生宕机 |
| FieldByIndex(index []int) StructField        | 多层成员访问时，根据 []int 提供的每个结构体的字段索引，返回字段的信息。没有找到时返回零值。当类型不是结构体或索引超界时 发生宕机 |
| FieldByNameFunc( match func(string) bool) (StructField,bool)        | 根据匹配函数匹配需要的字段。当值不是结构体或索引超界时发生宕机 |

#### 结构体字段类型

reflect.Type 的 Field() 方法返回 StructField 结构，这个结构描述结构体的成员信息，通过这个信息可以获取成员与结构体的关系，如偏移、索引、是否为匿名字段、结构体标签（Struct Tag）等，而且还可以通过 StructField 的 Type 字段进一步获取结构体成员的类型信息。StructField 的结构如下：

```go
// A StructField describes a single field in a struct.
type StructField struct {
	// Name is the field name.
	Name string
	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name. It is empty for upper case (exported) field names.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	PkgPath string

	Type      Type      // field type
	Tag       StructTag // field tag string
	Offset    uintptr   // offset within struct, in bytes
	Index     []int     // index sequence for Type.FieldByIndex
	Anonymous bool      // is an embedded field
}
```

字段说明如下。

- Name：为字段名称。
- PkgPath：字段在结构体中的路径。
- Type：字段本身的反射类型对象，类型为 reflect.Type，可以进一步获取字段的类型信息。
- Tag：结构体标签，为结构体字段标签的额外信息，可以单独提取。
- Index：FieldByIndex 中的索引顺序。
- Anonymous：表示该字段是否为匿名字段。

#### 获取成员反射信息

下面代码中，实例化一个结构体并遍历其结构体成员，再通过 reflect.Type 的 FieldByName() 方法查找结构体中指定名称的字段，直接获取其类型信息。

反射访问结构体成员类型及信息：

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 声明【猫】结构体
	type cat struct {
		Name string
		Type int `json:"type" id:"100"` // 带有标签字段
	}

	// 创建结构体实例
	instance := cat{Name: "mini", Type: 1}

	// 获取结构体实例的反射实例的反射类型对象
	typeOfCat := reflect.TypeOf(instance)

	// 遍历结构体所有成员
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取结构体字段类型
		field := typeOfCat.Field(i)

		// 输出成员名和标签
		fmt.Printf("name: %s, tag: '%s'\n", field.Name, field.Tag)
	}

	// 通过字段名，获取字段类型信息
	if field, exist := typeOfCat.FieldByName("Type"); exist {
		// 从标签中取出需要的具体标签
		fmt.Printf("json tag statement: %s, id tag statement: %s\n", field.Tag.Get("json"), field.Tag.Get("id"))
	}
}
```

代码输出如下：

```text
name: Name, tag: ''
name: Type, tag: 'json:"type" id:"100"'
json tag statement: type, id tag statement: 100
```

代码说明如下：

- 第 10 行，声明了带有两个成员的 cat 结构体。
- 第 12 行，Type 是 cat 的一个成员，这个成员类型后面带有一个以`开始和结尾的字符串。这个字符串在 Go 语言中被称为 Tag（标签）。一般用于给字段添加自定义信息，方便其他模块根据信息进行不同功能的处理。
- 第 16 行，创建 cat 实例，并对两个字段赋值。结构体标签属于类型信息，无须且不能赋值。
- 第 19 行，获取实例的反射类型对象。
- 第 22 行，使用 reflect.Type 类型的 NumField() 方法获得一个结构体类型共有多少个字段。如果类型不是结构体，将会触发宕机错误。
- 第 24 行，reflect.Type 中的 Field() 方法和 NumField 一般都是配对使用，用来实现结构体成员的遍历操作。
- 第 27 行，使用 reflect.Type 的 Field() 方法返回的结构不再是 reflect.Type 而是StructField 结构体。
- 第 31 行，使用 reflect.Type 的 FieldByName() 根据字段名查找结构体字段信息，cat Type 表示返回的结构体字段信息，类型为 StructField，ok 表示是否找到结构体字段的信息。
- 第 33 行中，使用 StructField 中 Tag 的 Get() 方法，根据 Tag 中的名字进行信息获取。