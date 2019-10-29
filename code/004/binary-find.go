package main

import (
	"flag"
	"fmt"
	"math"
)

func main() {
	// 定义命令行参数
	find := flag.Int("find", 0, "Please input integer value!")
	// 解析命令行参数
	flag.Parse()

	// 定义一个指定范围的数组(切片)
	s := generateSectionIntSlice(1, 50, 2)
	// 输出当前数组(切片)具体的长度
	fmt.Printf("slice: %+v, len: %d\n", s, len(s))

	// 待查找元素值假想索引值
	think := -1
	// 从数组中查找元素值，返回元素值索引值
	index := binaryFind(&s, 0, len(s)-1, *find, &think)

	// 输出查找结果
	fmt.Printf("Look for the element value %d, with the index value at position :%d in the slice\n", *find, index)
}

// 二分查找函数
// 假设有序数组的顺序是从小到大(这个很关键也很重要决定左右方向)
// arr 为待查找元素值源数组
// leftIndex 为源数组第一个元素索引值
// rightIndex 为源数组最后一个元素索引值
// think 为待查找元素假想索引值(-1为未查找到)
func binaryFind(arr *[]int, leftIndex, rightIndex int, findValue int, think *int) int {
	// 判断leftIndex是否大于rightIndex，如果大于了，则直接返回-1就OK
	if leftIndex > rightIndex {
		*think = -1
		// 未找到
		return *think
	}

	// 先找到中间元素的下标，这个地方一定要特别注意，只要涉及到除法运行肯定会出现小数，则这里要特殊处理一下
	middle := int(math.Ceil(float64((leftIndex + rightIndex) / 2)))

	// 二分查找一共分为四种情况(这个分支语句一共处理三种情况)
	if (*arr)[middle] > findValue {
		// 在左侧查找(查找范围应该在 leftIndex ~ middle - 1的位置查找)
		binaryFind(arr, leftIndex, middle-1, findValue, think)
	} else if (*arr)[middle] < findValue {
		// 在右侧查找(查找范围应该在 middle+1 ~ rightIndex的位置查找)
		binaryFind(arr, middle+1, rightIndex, findValue, think)
	} else {
		*think = middle
		// 查找的元素值正好是中间索引位置
		return *think
	}

	// 未找到
	return *think
}

// 生成指定范围内的整型切片(数组)
func generateSectionIntSlice(min, max int, step int) []int {
	result := make([]int, 0, max)
	for i := min; i <= max; i += step {
		result = append(result, i)
	}
	return result
}
