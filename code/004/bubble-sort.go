package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	origin := generateSectionIntSliceOfDisorderly(1, 50)
	fmt.Printf("Origin Slice Containter: %+v, len(%d)\n", origin, len(origin))
	bubbleSort(origin)
	fmt.Printf("Bubble Sort Slice Containter: %+v, len(%d)\n", origin, len(origin))
}

// 生成指定范围的随机数
func Random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// 生成指定范围内的整型切片
func generateSectionIntSliceOfDisorderly(min, max int) []int {
	l := max - min
	r := make([]int, 0)
	for i := 0; i < l; i++ {
		r = append(r, Random(min, max))
	}
	return r
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
