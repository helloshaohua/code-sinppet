package main

import "fmt"

func main() {
	ch := make(chan int, 2)

	ch <- 1
	ch <- 2

	close(ch)

	for i := 0; i < cap(ch)+2; i++ {
		value, ok := <-ch
		fmt.Printf("value %d, ok: %t\n", value, ok)
	}
}
