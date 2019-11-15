package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	u := User{1, "张三", 20}
	va := reflect.ValueOf(u)
	vb := reflect.ValueOf(&u)
	// 值类型是不可以修改的
	fmt.Println(va.CanSet(), va.FieldByName("Name").CanSet())
	// 值类型是不可以修改的
	fmt.Println(vb.CanSet(), vb.Elem().FieldByName("Name").CanSet())
	vb.Elem().FieldByName("Name").SetString("zhangsan")
	fmt.Printf("%v\n", vb)
}
