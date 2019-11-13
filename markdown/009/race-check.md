### Go语言竞态检测(检测代码在并发环境下可能出现的问题{通过使用原子访问解决})


Go语言程序可以使用通道进行多个 goroutine 间的数据交换，但这仅仅是数据同步中的一种方法。通道内部的实现依然使用了各种锁，因此优雅代码的代价是性能。在某些轻量级的场合，原子访问（atomic包）、互斥锁（sync.Mutex）以及等待组（sync.WaitGroup）能最大程度满足需求。

本节只讲解原子访问，互斥锁和等待组将在接下来的两节中讲解。

当多线程并发运行的程序竞争访问和修改同一块资源时，会发生竞态问题。

下面的代码中有一个 ID 生成器，每次调用生成器将会生成一个不会重复的顺序序号，使用 10 个并发生成序号，观察 10 个并发后的结果。

> 竞态检测的具体代码：

```go
package main

import (
	"fmt"
	"sync/atomic"
)

// 序列号存储变量初始化
var code int64

// 序列号生成器
func GenerateID() int64 {
	// 序列号存储值加1
	atomic.AddInt64(&code, 1)

	// 返回序列号
	return code
}

func main() {
	// 生成10个并发序列号
	for i := 0; i < 10; i++ {
		go GenerateID()
	}

	// 输出并发序列号
	fmt.Println(GenerateID())
}
```

> 代码说明如下：

- 第 9 行，序列号生成器中的保存上次序列号的变量。
- 第 14 行，使用原子操作函数 atomic.AddInt64() 对 code 变量加 1 操作。不过这里故意没有使用 atomic.AddInt64() 的返回值作为 GenerateID() 函数的返回值，因此会造成一个竞态问题。
- 第 23 行，循环 10 次生成 10 个 goroutine 调用 GenerateID() 函数，同时忽略 GenerateID() 的返回值。
- 第 27 行，单独调用一次 GenerateID() 函数。

在运行程序时，为运行参数加入-race参数，开启运行时（runtime）对竞态问题的分析，命令如下：

```go
go run -race test.go
```

代码运行发生宕机，输出信息如下：

```text

==================
WARNING: DATA RACE
Write at 0x000001227320 by goroutine 8:
  sync/atomic.AddInt64()
      /usr/local/Cellar/go/1.13/libexec/src/runtime/race_amd64.s:276 +0xb
  main.GenerateID()
      /Users/warnerwu/dev/go/src/code-snippet/test.go:14 +0x43

Previous read at 0x000001227320 by goroutine 7:
  main.GenerateID()
      /Users/warnerwu/dev/go/src/code-snippet/test.go:17 +0x53

Goroutine 8 (running) created at:
  main.main()
      /Users/warnerwu/dev/go/src/code-snippet/test.go:23 +0x4f

Goroutine 7 (finished) created at:
  main.main()
      /Users/warnerwu/dev/go/src/code-snippet/test.go:23 +0x4f
==================
10
Found 1 data race(s)
exit status 66
```

根据报错信息，第 17 行有竞态问题，根据 atomic.AddInt64() 的参数声明，这个函数会将修改后的值以返回值方式传出。对 GenerateID 函数修改如下：

```go
// 序列号生成器
func GenerateID() int64 {
	// 返回序列号
	return atomic.AddInt64(&code, 1)
}
```

再次运行：

```go
go run -race test.go
```

代码输出如下：

```text
11
```
    
没有发生竞态问题，程序运行正常。

本例中只是对变量进行增减操作，虽然可以使用互斥锁（sync.Mutex）解决竞态问题，但是对性能消耗较大。在这种情况下，推荐使用原子操作（atomic）进行变量操作。





