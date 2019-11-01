### Go语言冒泡排序

冒泡排序算法（bubble sort）是一种很简单的交换排序，每轮都从第一个元素开始，依次将较大值向后交换一位，直至整个队列排序完成。

示例代码如下所示：

```go
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

```

执行结果如下所示：

```text
Origin Slice Containter: [29 20 15 12 9 1 29 41 4 35 33 30 29 47 1 24 18 43 34 23 23 8 23 25 2 4 43 1 28 25 41 11 37 41 2 32 33 25 2 16 4 43 45 42 18 20 5 30 1], len(49)
Bubble Sort Slice Containter: [1 1 1 1 2 2 2 4 4 4 5 8 9 11 12 15 16 18 18 20 20 23 23 23 24 25 25 25 28 29 29 29 30 30 32 33 33 34 35 37 41 41 41 42 43 43 43 45 47], len(49)
```
