### Go语言互斥锁(sync.Mutex)和读写互斥锁(sync.RWMutex)

Go语言包中的 sync 包提供了两种锁类型：sync.Mutex 和 sync.RWMutex。

Mutex 是最简单的一种锁类型，同时也比较暴力，当一个 goroutine 获得了 Mutex 后，其他 goroutine 就只能乖乖等到这个 goroutine 释放该 Mutex。

RWMutex 相对友好些，是经典的单写多读模型。在读锁占用的情况下，会阻止写，但不阻止读，也就是多个 goroutine 可同时获取读锁（调用 RLock() 方法；而写锁（调用 Lock() 方法）会阻止任何其他 goroutine（无论读和写）进来，整个锁相当于由该 goroutine 独占。从 RWMutex 的实现看，RWMutex 类型其实组合了 Mutex：

```go
// A RWMutex is a reader/writer mutual exclusion lock.
single-address
// The zero value for a RWMutex is an unlocked mutex.
//
// A RWMutex must not be copied after first use.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
type RWMutex struct {
	w           Mutex  // held if there are pending writers
	writerSem   uint32 // semaphore for writers to wait for completing readers
	readerSem   uint32 // semaphore for readers to wait for completing writers
	readerCount int32  // number of pending readers
	readerWait  int32  // number of departing readers
}
```

对于这两种锁类型，任何一个 Lock() 或 RLock() 均需要保证对应有 Unlock() 或 RUnlock() 调用与之对应，否则可能导致等待该锁的所有 goroutine 处于饥饿状态，甚至可能导致死锁。

#### 互斥锁

锁的典型使用模式如下：

```go
package main

import (
	"fmt"
	"sync"
)

// 变量列表声明
var (
	count      int        // 逻辑中使用的计数变量
	countGuard sync.Mutex // 与变量对应的互斥锁使用
)

// GetCount 获取计数器
func GetCount() int {
	// 锁定
	countGuard.Lock()

	// 在函数退出时解除锁定
	defer countGuard.Unlock()

	// 返回计数
	return count
}

// SetCount 设置计数器
func SetCount(number int) {
	// 锁定
	countGuard.Lock()

	// 改变计数器值
	count = number

	// 解除锁定
	countGuard.Unlock()
}

func main() {
	// 可以进行并发安全的设置计数器值
	SetCount(1)

	// 可以进行并发安全的获取计数器值
	fmt.Println(GetCount())
}
```

代码说明如下：

- 第 10 行，是逻辑中使用的计数变量，无论是包级的变量还是结构体成员字段，都可以。
- 第 11 行，一般情况下，建议将互斥锁的粒度设置得越小越好，降低因为共享访问时等待的时间。这里笔者习惯性地将互斥锁的变量命名为以下格式：`变量名+Guard` 以表示这个互斥锁用于保护这个变量。
- 第 15 行，是一个获取 count 值的函数封装，通过这个函数可以并发安全的访问变量 count。
- 第 17 行，尝试对 countGuard 互斥量进行加锁。一旦 countGuard 发生加锁，如果另外一个 goroutine 尝试继续加锁时将会发生阻塞，直到这个 countGuard 被解锁。
- 第 20 行，使用 defer 将 countGuard 的解锁进行延迟调用，解锁操作将会发生在 GetCount() 函数返回时。
- 第 27 行，在设置 count 值时，同样使用 countGuard 进行加锁、解锁操作，保证修改 count 值的过程是一个原子过程，不会发生并发访问冲突。

#### 读写互斥锁

在读多写少的环境中，可以优先使用读写互斥锁（sync.RWMutex），它比互斥锁更加高效。sync 包中的 RWMutex 提供了读写互斥锁的封装。

我们将互斥锁例子中的一部分代码修改为读写互斥锁，参见下面代码：

```go
package main

import (
	"fmt"
	"sync"
)

// 变量列表声明
var (
	count      int          // 逻辑中使用的计数变量
	countGuard sync.RWMutex // 与变量对应的读写互斥锁使用
)

// GetCount 获取计数器
func GetCount() int {
	// 使用读写互斥锁锁定
	countGuard.RLock()

	// 在函数退出时解除锁定
	defer countGuard.RUnlock()

	// 返回计数
	return count
}

// SetCount 设置计数器
func SetCount(number int) {
	// 改变计数器值
	count = number
}

func main() {
	// 可以进行并发安全的设置计数器值
	SetCount(1)

	// 可以进行并发安全的获取计数器值
	fmt.Println(GetCount())
}
```

代码说明如下：

- 第 11 行，在声明 countGuard 时，从 sync.Mutex 互斥锁改为 sync.RWMutex 读写互斥锁。
- 第 17 行，获取 count 的过程是一个读取 count 数据的过程，适用于读写互斥锁。在这一行，把 countGuard.Lock() 换做 countGuard.RLock()，将读写互斥锁标记为读状态。如果此时另外一个 goroutine 并发访问了 countGuard，同时也调用了 countGuard.RLock() 时，并不会发生阻塞。
- 第 20 行，与读模式加锁对应的，使用读模式解锁。
