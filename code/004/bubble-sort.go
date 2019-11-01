package main

import (
	"fmt"

	"github.com/wumoxi/toolkit"
)

func main() {
	origin := toolkit.GenerateSectionIntSliceOfDisorderly(1, 20)
	fmt.Printf("Origin Slice Containter: %+v, len(%d)\n", origin, len(origin))
	bubbleSort(origin)
	fmt.Printf("Bubble Sort Slice Containter: %+v, len(%d)\n", origin, len(origin))
}

// 冒泡排序
func bubbleSort(slice []int) []int {
	last := len(slice)
	for i := 0; i < last; i++ {
		for j := i; j < last; j++ {
			if slice[i] > slice[j] {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
	return slice
}
