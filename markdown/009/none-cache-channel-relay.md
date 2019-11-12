### Go语言无缓冲的通道模拟接力赛

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// 用 WaitGroup 等待Goroutine程序结束
var wg sync.WaitGroup

func main() {
	// 创建一个无缓冲的通道用于传递接力棒
	relay := make(chan int)

	// 为最后一位跑步者，添加Goroutine等待计数
	wg.Add(1)

	// 第一位跑步者持有接力棒
	go Runner(relay)

	// 开始比赛
	relay <- 1

	// 等待比赛结束
	wg.Wait()
}

// Runner 模拟接力赛中的每一位跑步者
func Runner(replay chan int) {
	// 表示每一位新的跑步者
	var newRunner int

	// 等待接力棒
	runner := <-replay

	// 开始围绕跑道跑步
	fmt.Printf("Runner %d runner with relay\n", runner)

	// 判断当前跑步者是不是第四位，不是则创建下一位跑步者
	if runner != 4 {
		newRunner = runner + 1
		fmt.Printf("Runner %d to the line\n", newRunner)
		go Runner(replay)
	}

	// 围绕跑道跑步
	time.Sleep(100 * time.Millisecond)

	// 比赛结束了吗？
	if runner == 4 {
		fmt.Printf("Runner %d finished, race over\n", runner)
		wg.Done()
		return
	}

	// 传递接力棒
	replay <- newRunner
}
```

运行这个程序，输出结果如下所示。

```text
Runner 1 runner with relay
Runner 2 to the line
Runner 2 runner with relay
Runner 3 to the line
Runner 3 runner with relay
Runner 4 to the line
Runner 4 runner with relay
Runner 4 finished, race over
```

在这两个例子里，我们使用无缓冲的通道同步 goroutine，模拟了网球和接力赛。代码的流程与这两个活动在真实世界中的流程完全一样，这样的代码很容易读懂。基本上现在知道了无缓冲的通道是如何工作的了。