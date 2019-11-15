### 对结构的反射操作获取结构体成员信息

下面的示例演示了如何获取一个结构中所有成员的值：

```go
package main

import (
	"fmt"
	"reflect"
)

// Student 学生结构体类型
type Student struct {
	Name   string  // 姓名
	Age    int     // 年龄
	Sex    uint8   // 1 男， 2 女
	Weight float64 // 体重
}

func main() {
	// 用户实例
	student := Student{Name: "张三", Age: 18, Sex: 1, Weight: 68.57}

	// 获取实例的值
	v := reflect.ValueOf(&student).Elem()

	// 获取实例的值的类型
	t := v.Type()

	// 通过【实例的值】的方法【NumField】获取字段的个数
	for i := 0; i < v.NumField(); i++ {
		// 通过【实例的值】的方法【Field】获取字段
		field := v.Field(i)

		// 字段名称
		fieldName := t.Field(i).Name

		// 字段类型
		fieldType := field.Type()

		// 字段值
		fieldValue := field.Interface()

		// 输出结构体成员信息
		fmt.Printf("字段索引值：%d, 字段名称：%s, 字段类型：%s, 字段值：%v\n", i, fieldName, fieldType, fieldValue)
	}
}
```

以上例子的输出为：

```text
字段索引值：0, 字段名称：Name, 字段类型：string, 字段值：张三
字段索引值：1, 字段名称：Age, 字段类型：int, 字段值：18
字段索引值：2, 字段名称：Sex, 字段类型：uint8, 字段值：1
字段索引值：3, 字段名称：Weight, 字段类型：float64, 字段值：68.57
```

可以看出，对于结构的反射操作并没有根本上的不同，只是用了 Field() 方法来按索引获取对应的成员。