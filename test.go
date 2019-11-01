package main

import (
	"fmt"
	"toolkit"
)

func sum(input ...int) (sum int) {
	for _, value := range input {
		sum += value
	}
	return sum
}

func main() {
	fmt.Println(sum(toolkit.GenerateSectionIntSliceOfDisorderly(1, 100)...))
}
