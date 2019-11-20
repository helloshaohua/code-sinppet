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
