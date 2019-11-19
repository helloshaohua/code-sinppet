package main

import (
	"fmt"
	"reflect"
)

func main() {
	x := 2                   // value type variable?
	a := reflect.ValueOf(2)  // 2 int no
	b := reflect.ValueOf(x)  // 2 int no
	c := reflect.ValueOf(&x) // &x *int no
	d := c.Elem()            // 2 int yes

	fmt.Println(a.CanAddr()) // false
	fmt.Println(b.CanAddr()) // false
	fmt.Println(c.CanAddr()) // false
	fmt.Println(d.CanAddr()) // true

}
