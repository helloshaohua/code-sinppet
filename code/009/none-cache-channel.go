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
