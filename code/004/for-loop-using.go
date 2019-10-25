package main

import "fmt"

func main() {
	level := 9
	for i := 1; i <= level; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%d * %d = %d\t", i, j, i*j)
		}
		fmt.Println()
	}
}
