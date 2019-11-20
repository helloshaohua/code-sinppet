### Go语言通过类型信息创建实例

当已知 reflect.Type 时，可以动态地创建这个类型的实例，实例的类型为指针。例如 reflect.Type 的类型为 int 时，创建 int 的指针，即*int，代码如下：

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 声明整型变量 a
	var a int

	// 获取变量 a 的反射类型对象
	typeOfA := reflect.TypeOf(a)

	// 根据反射类型对象创建类型实例
	instance := reflect.New(typeOfA)

	// 输出Value的类型和种类
	fmt.Println(instance.Type(), instance.Kind())

	// 修改类型实例值
	instance.Elem().SetInt(2048)

	// 输出类型值
	fmt.Printf("instance: %+v\n", instance.Elem().Int())
}
```

代码输出如下：

```text
*int ptr
instance: 2048
```

代码说明如下：

- 第 13 行，获取变量 a 的反射类型对象。
- 第 16 行，使用 reflect.New() 函数传入变量 a 的反射类型对象，创建这个类型的实例值，值以 reflect.Value 类型返回。这步操作等效于：new(int)，因此返回的是 *int 类型的实例。
- 第 19 行，打印 instance 的类型为 *int，种类为指针。
- 第 22 行，修改 instance 类型的值为 2048。
- 第 25 行，输出 instance 类型的值。
