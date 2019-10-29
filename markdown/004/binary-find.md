### 

二分查找是 logN 级别的查找算法，前提是有序序列并且存储在顺序表中，如果存储方式是链表则不能使用。

二分查找的核心思想理解起来非常简单，有点类似分治思想，即每次都通过跟区间中的中间元素对比，将待查找的区间缩小为一半，直到找到要查找的元素，或者区间被缩小为 0，但是二分查找的代码实现比较容易写错，需要着重掌握它的三个容易出错的地方：循环退出条件、mid 的取值，left 和 right 的更新。

二分查找虽然性能比较优秀，但应用场景也比较有限，底层必须依赖数组，并且还要求数据是有序的，对于较小规模的数据查找，我们直接使用顺序遍历就可以了，二分查找的优势并不明显，二分查找更适合处理静态数据，也就是没有频繁插入、删除操作的数据。
思路：
1) 先找到中间的下标 middle = (leftIndex + RightIndex) /2 ，然后用中间的下标值和 FindVal 比较。

```text
a：如果 arr[middle] > FindVal，那么就向 LeftIndex ~ (midlle-1) 区间找。
b：如果 arr[middle] < FindVal，那么就向 middle+1 ~ RightIndex 区间找。
c：如果 arr[middle] == FindVal，那么直接返回。
```

2) 从第一步的 a、b、c 递归执行，直到找到位置。

3) 如果 LeftIndex > RightIndex，则表示找不到，退出。


代码/举例：
定义一个包含（1, 2, 5, 7, 15, 25, 30, 36, 39, 51, 67, 78, 80, 82, 85, 91, 92, 97）等值的数组，假设说要查找 30 这个值，如果按照循环的查找方法，找到 30 这个值要执行 7 次，那么如果是按照二分查找呢？二分查找的过程如下：

```text
left = 1, right = 18; mid = (1+18)/2 = 9; 51 > 30
left = 1, right = mid - 1 = 8; mid = (1+8)/2 = 4; 15 < 30
left = mid + 1 = 5, right = 8; mid = (5+8)/2 = 6; 30 = 30 查找完毕
```

如上所示只需要执行 3 次，大大减少了执行时间，具体代码实现如下所示：

```go
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
```

执行结果如下所示：

```shell
go run binary-find.go -find 23
slice: [1 3 5 7 9 11 13 15 17 19 21 23 25 27 29 31 33 35 37 39 41 43 45 47 49], len: 25
Look for the element value 23, with the index value at position :11 in the slice
```