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
	vb := reflect.ValueOf(&u)
	fmt.Printf("%v\n", vb)
}
