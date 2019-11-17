package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Print 打印值x的方法集
func Print(x interface{}) {
	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)
	fmt.Printf("type %s\n", t)
	for i := 0; i < v.NumMethod(); i++ {
		mt := v.Method(i).Type()
		fmt.Printf("func (%s) %s %s\n", t, t.Method(i).Name, strings.TrimPrefix(mt.String(), "func"))
	}
	fmt.Println()
}

func main() {
	Print(time.Hour)
	Print(new(strings.Replacer))
}
