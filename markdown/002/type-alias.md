### 在结构体成员嵌入时使用别名

当类型别名作为结构体嵌入的成员时会发生什么情况呢？请参考下面的代码。

```go
package main

import (
	"fmt"
	"reflect"
)

// 定义商标结构
type Brand struct {
}

// 为商标结构添加Show方法
func (t Brand) Show() {
}

// 为Brand定义一个别名FakeBrand
type FakeBrand = Brand

// 定义车辆结构
type Vehicle struct {
	// 嵌入两个结构
	FakeBrand
	Brand
}

func main() {
	// 声明变量a为车辆类型
	var a Vehicle

	// 指定调用FakeBrand的Show
	a.FakeBrand.Show()

	// 获取a的类型反射对象
	ta := reflect.TypeOf(a)

	// 遍历a的所有成员
	for i := 0; i < ta.NumField(); i++ {
		// a的成员信息
		field := ta.Field(i)

		// 打印成员的字段名和类型
		fmt.Printf("FieldName: %v, FieldType: %v\n", field.Name, field.Type.Name())
	}
}
```
代码输出如下：
```text
FieldName: FakeBrand, FieldType: Brand
FieldName: Brand, FieldType: Brand
```

代码说明如下：

- 第 9 行，定义商标结构。
- 第 13 行，为商标结构添加 Show() 方法。
- 第 17 行，为 Brand 定义一个别名 FakeBrand。
- 第 20～24 行，定义车辆结构 Vehicle，嵌入 FakeBrand 和 Brand 结构。
- 第 28 行，将 Vechicle 实例化为 a。
- 第 31 行，显式调用 Vehicle 中 FakeBrand 的 Show() 方法。
- 第 34 行，使用反射取变量 a 的反射类型对象，以查看其成员类型。
- 第 37～43 行，遍历 a 的结构体成员。
- 第 42 行，打印 Vehicle 类型所有成员的信息。

这个例子中，FakeBrand 是 Brand 的一个别名，在 Vehicle 中嵌入 FakeBrand 和 Brand 并不意味着嵌入两个 Brand，FakeBrand 的类型会以名字的方式保留在 Vehicle 的成员中。

如果尝试将第 31 行改为：

```go
a.Show()
```

编译器将发生报错：

```shell
ambiguous selector a.Show
```

在调用 Show() 方法时，因为两个类型都有 Show() 方法，会发生歧义，证明 FakeBrand 的本质确实是 Brand 类型。


