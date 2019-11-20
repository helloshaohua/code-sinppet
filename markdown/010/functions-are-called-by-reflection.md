### Go语言通过反射调用函数

如果反射值对象（reflect.Value）中值的类型为函数时，可以通过 reflect.Value 调用该函数。使用反射调用函数时，需要将参数使用反射值对象的切片 []reflect.Value 构造后传入 Call() 方法中，调用完成时，函数的返回值通过 []reflect.Value 返回。

下面的代码声明一个加法函数，传入两个整型值，返回两个整型值的和。将函数保存到反射值对象（reflect.Value）中，然后将两个整型值构造为反射值对象的切片（[]reflect.Value），使用 Call() 方法进行调用。

反射调用函数：

```go
package main

import (
	"fmt"
	"reflect"
)

// 学生结构体类型
type Student struct {
	Username string
	Age      int
	Sex      int
	Mobile   string
}

// 普通函数获取用户实例信息
func getUserInfo(username string) Student {
	return Student{"武沫汐", 18, 1, "13729282822"}
}

// 普通函数
func add(a, b int) int {
	return a + b
}

func main() {
	// 将函数包装为反射值对象
	funcOfValue := reflect.ValueOf(add)

	// 构造函数参数，转入两个整型值
	paramsList := []reflect.Value{reflect.ValueOf(15), reflect.ValueOf(25)}

	// 通过反射调用函数
	resultList := funcOfValue.Call(paramsList)

	// 获取第一个返回值，取整数值
	fmt.Println(resultList[0].Int())

	// 将函数包装为返回值对象
	fOfV := reflect.ValueOf(getUserInfo)

	// 构造函数参数，转入两个整型值
	paramsList = []reflect.Value{reflect.ValueOf("武沫汐")}

	// 通过反射调用函数
	resultList = fOfV.Call(paramsList)

	// 获取第一个返回值，取整数值
	fmt.Println(resultList[0])
}
```

代码说明如下：

- 第 9～14 行，定义学习结构体类型。
- 第 17～19 行，定义一个获取用户实例信息函数。
- 第 22～24 行，定义一个普通的加法函数。
- 第 28 行，将 add 函数包装为反射值对象。
- 第 31 行，将 10 和 20 两个整型值使用 reflect.ValueOf 包装为 reflect.Value，再将反射值对象的切片 []reflect.Value 作为函数的参数。
- 第 34 行，使用 funcValue 函数值对象的 Call() 方法，传入参数列表 paramList 调用 add() 函数。
- 第 37 行，调用成功后，通过 retList[0] 取返回值的第一个参数，使用 Int 取返回值的整数值。

**提示**

反射调用函数的过程需要构造大量的 reflect.Value 和中间变量，对函数参数值进行逐一检查，还需要将调用参数复制到调用函数的参数内存中。调用完毕后，还需要将返回值转换为 reflect.Value，用户还需要从中取出调用值。因此，反射调用函数的性能问题尤为突出，不建议大量使用反射函数调用。