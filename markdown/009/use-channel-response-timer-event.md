### Go语言使用通道响应计时器的事件

Go语言中的 time 包提供了计时器的封装。由于 Go语言中的通道和 goroutine 的设计，定时任务可以在 goroutine 中通过同步的方式完成，也可以通过在 goroutine 中异步回调完成。这里将分两种用法进行例子展示。

#### 一段时间之后（time.After）

延迟回调：

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// 声明一个退出用的通道
	exit := make(chan bool)

	// 打印开始
	fmt.Println("start")

	// 过1秒钟后，调用匿名函数
	time.AfterFunc(time.Second, func() {
		// 1秒钟后，打印结束
		fmt.Println("one second after-func")

		// 通知main函数的 goroutine 已经结束
		exit <- true
	})

	// 等待结束
	<-exit
}
```

代码说明如下：
- 第 10 行，声明一个退出用的通道，往这个通道里写数据表示退出。
- 第 16 行，调用 time.AfterFunc() 函数，传入等待的时间和一个回调。回调使用一个匿名函数，在时间到达后，匿名函数会在另外一个 goroutine 中被调用。
- 第 21 行，任务完成后，往退出通道中写入数值表示需要退出。
- 第 25 行，运行到此处时持续阻塞，直到 1 秒后第 21 行被执行后结束阻塞。

time.AfterFunc() 函数是在 time.After 基础上增加了到时的回调，方便使用。

而 time.After() 函数又是在 time.NewTimer() 函数上进行的封装，下面的例子展示如何使用 timer.NewTimer() 和 time.NewTicker()。

#### 定点计时

计时器（Timer）的原理和倒计时闹钟类似，都是给定多少时间后触发。打点器（Ticker）的原理和钟表类似，钟表每到整点就会触发。这两种方法创建后会返回 time.Ticker 对象和 time.Timer 对象，里面通过一个 C 成员，类型是只能接收的时间通道（<-chan Time），使用这个通道就可以获得时间触发的通知。

下面代码创建一个打点器，每 500 毫秒触发一起；创建一个计时器，2 秒后触发，只触发一次。

计时器：

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// 创建一个打点器，每500毫秒触发一次
	ticker := time.NewTicker(time.Millisecond * 500)

	// 创建一个计时器，2秒后触发
	stopper := time.NewTimer(time.Second * 2)

	// 声明计数变量
	var count int

	// 不断地检查通道情况
	for {
		// 多路复用通道
		select {
		case <-stopper.C: // 计时器到时了
			fmt.Println("stop")
			// 跳出循环
			goto stopHere
		case <-ticker.C: // 打点器触发了
			// 记录触发了多少次
			count++
			fmt.Printf("tick number %d\n", count)
		}
	}

// 退出的标签，使用goto语句跳转
stopHere:
	fmt.Println("Done")
}
```

代码说明如下：
- 第 10 行，创建一个打点器，500 毫秒触发一次，返回 *time.Ticker 类型变量。
- 第 13 行，创建一个计时器，2 秒后返回，返回 *time.Timer 类型变量。
- 第 16 行，声明一个变量，用于累计打点器触发次数。
- 第 19 行，每次触发后，select 会结束，需要使用循环再次从打点器返回的通道中获取触发通知。
- 第 21 行，同时等待多路计时器信号。
- 第 22 行，计时器信号到了。
- 第 25 行，通过 goto 跳出循环。
- 第 26 行，打点器信号到了，通过i自加记录触发次数并打印。