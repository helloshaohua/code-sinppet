package main

import (
	"fmt"
	"time"
)

func main() {
	input := make(chan interface{})

	// product - product the message
	go func() {
		for i := 0; i < 5; i++ {
			input <- i
		}
		input <- "hello world"
	}()

	t1 := time.NewTimer(5 * time.Second)
	t2 := time.NewTimer(10 * time.Second)

	for {
		select {
		case message := <-input:
			fmt.Println(message)
		case <-t1.C:
			fmt.Println("5s second")
			t1.Reset(5 * time.Second)
		case <-t2.C:
			fmt.Println("10s second")
			t2.Reset(10 * time.Second)
		}
	}
}
