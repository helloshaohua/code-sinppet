package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.4
	reflect.ValueOf(&x).Elem().SetFloat(5.5)
	fmt.Println(reflect.ValueOf(&x).Elem())
	fmt.Println(x)
}
