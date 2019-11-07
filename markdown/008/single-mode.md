### Go语言单例模式

设计模式的重要性不用多说，也是面试时常常会被问到的问题。但是对于设计模式，更多的则是仁者见仁智者见智。要在实际工作中不断的积累，再进行深度的思考，才能逐渐形成的一种思维。

单例模式就是一种设计模式，它能保证系统运行中一个类只创建一个实例。设想一下，系统中某个类用来加载配置文件，如果每次加载都创建一个实例，是不是就会造成资源浪费呢？这时候使用单例模式就可以节省多次加载的内存消耗。

#### 单例模式实现

单例模式可以分为懒汉式和饿汉式。懒汉式就是创建对象时比较懒，先不急着创建对象，在需要加载配置文件的时候再去创建。饿汉式就是在系统初始化的时候我们已经把对象创建好了，需要用的时候直接拿过来用就好了。不管那种模式最终目的只有一个，就是只实例化一次，只允许一个实例存在。

但是，在Go语言的世界中，没有 private、public、static 等关键字，也没有面向对象中类的概念。那么Go语言是如何控制访问范围的呢？首字母大写，代表对外部可见，首字母小写代表对外部不可见，适用于所有对象，包括函数、方法。

#### 懒汉式

```go
package main

type config struct {
}

var cfg *config

func GetConfigInstance() *config {
	if cfg == nil {
		cfg = new(config)
		return cfg
	}
	return cfg
}
```

上述代码没有考虑线程安全，我们可以使用Go语言 sync.Once 结构体，它提供了一个 Do 方法，该方法只在第一次调用时执行，从而保证了线程的安全，实现单例模式。

```go
package main

import "sync"

type config struct {
}

var cfg *config
var once sync.Once

func GetConfigInstance() *config {
	once.Do(func() {
		cfg = new(config)
	})
	return cfg
}
```

#### 饿汉式

Go语言饿汉式可以使用 init 函数，也可以使用全局变量。

```go
package main

type config struct {
}

var cfg *config

func init() {
	cfg = new(config)
}

// GetConfigInstance 提供获取实例的方法
func GetConfigInstance() *config {
	return cfg
}
```

```go
package main

type config struct {
}

var cfg *config = new(config)

// GetConfigInstance 提供获取实例的方法
func GetConfigInstance() *config {
	return cfg
}
```