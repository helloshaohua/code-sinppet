// 这个示例程序展示如何用无缓冲的通道来模拟，2 个goroutine 间的网球比赛
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ch := make(chan int)

	wg.Add(2)

	go player("张三", ch)
	go player("李四", ch)

	ch <- 1

	wg.Wait()
}

func player(name string, ch chan int) {
	// 函数退出时减少一个goroutine 等待计数
	defer wg.Done()
	for {
		ball, ok := <-ch
		// 通道已经关闭表示对方已经丢球了，我们赢了
		if !ok {
			fmt.Printf("Player %s Won\n", name)
			return
		}

		// 表示我们丢球了，我们输了
		if rand.Intn(100)%15 == 0 {
			fmt.Printf("Player %s Missed\n", name)
			close(ch)
			return
		}

		fmt.Printf("Player %s Hit %d\n", name, ball)
		ball++

		ch <- ball
	}
}
