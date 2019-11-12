// 这个示例程序展示如何用无缓冲的通道来模拟，2 个goroutine 间的网球比赛
package main

import (
	"fmt"
	"sync"
	"time"
)

// 用 WaitGroup 等待Goroutine程序结束
var wg sync.WaitGroup

// main 是所有Go程序的入口
func main() {
	// 创建一个无缓冲的通道
	relay := make(chan int)

	// 为最后一个跑步者将计数加1
	wg.Add(1)

	// 第一位跑步者持有接力棒
	go Runner(relay)

	// 开始比赛
	relay <- 1

	// 等待比赛结束
	wg.Wait()
}

// runner 模拟接力赛中的一位跑步者
func Runner(relay chan int) {
	var newRunner int

	// 等待接力棒
	runner := <-relay

	// 开始绕着跑道跑步
	fmt.Printf("Runner %d running with relay\n", runner)

	// 创建下一位跑步者
	if runner != 4 {
		newRunner = runner + 1
		fmt.Printf("Runner %d to the line\n", newRunner)
		go Runner(relay)
	}

	// 围绕跑道跑
	time.Sleep(100 * time.Millisecond)

	// 比赛结束了吗？
	if runner == 4 {
		fmt.Printf("Runner %d finished, race over\n", runner)
		wg.Done()
		return
	}

	// 将接力棒交给下一们跑步者
	fmt.Printf("Runner %d exchangeed with runner %d\n", runner, newRunner)

	relay <- newRunner
}
