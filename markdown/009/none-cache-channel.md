### Go语言无缓冲的通道模拟网球比赛

在网球比赛中，两位选手会把球在两个人之间来回传递。选手总是处在以下两种状态之一：要么在等待接球，要么将球打向对方。可以使用两个 goroutine 来模拟网球比赛，并使用无缓冲的通道来模拟球的来回，代码如下所示。

```go
// 这个示例程序展示如何用无缓冲的通道来模拟，2 个goroutine 间的网球比赛
package _09

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 用来等待程序结束
var wg sync.WaitGroup

// 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// 创建一个无缓冲的通道
	ch := make(chan int)

	// 添加等待计数器为2，表示等待两个 goroutine
	wg.Add(2)

	// 启动两个选手
	go player("张三", ch)
	go player("李四", ch)

	// 发球
	ch <- 1

	// 等待游戏结束
	wg.Wait()
}

// player 模拟一个选手在打网球
func player(name string, ch chan int) {
	// 在函数退出时调用Done方法来通知main 函数工作已经完成
	defer wg.Done()

	for {
		// 等待网球被击打过来
		ball, ok := <-ch
		if !ok {
			// 如果通道被关闭，我们就赢了
			fmt.Printf("Player %s Won\n", name)
			return
		}

		// 选随机数，然后用这个随机数来判断我们是否丢球
		if rand.Intn(100)%13 == 0 {
			fmt.Printf("Player %s Missed\n", name)
			// 关闭通道，表示我们输了
			close(ch)
			return
		}

		// 显示击球数，并将击球数加1
		fmt.Printf("Player %s Hit %d\n", name, ball)
		ball++

		// 将网球打向对手
		ch <- ball
	}
}
```

运行这个程序，输出结果如下所示。

```text
Player 李四 Hit 1
Player 张三 Hit 2
Player 李四 Hit 3
Player 张三 Hit 4
Player 李四 Hit 5
Player 张三 Hit 6
Player 李四 Hit 7
Player 张三 Hit 8
Player 李四 Hit 9
Player 张三 Hit 10
Player 李四 Hit 11
Player 张三 Hit 12
Player 李四 Hit 13
Player 张三 Hit 14
Player 李四 Hit 15
Player 张三 Missed
Player 李四 Won
```