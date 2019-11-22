package main

import "fmt"

func main() {
	fmt.Printf("%+v, %T\n", 2&0x0001, 2&0x0001)
}
