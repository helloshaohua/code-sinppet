### Go语言使用事件系统实现事件的响应和处理

Go语言可以将类型的方法与普通函数视为一个概念，从而简化方法和函数混合作为回调类型时的复杂性。这个特性和 C# 中的代理（delegate）类似，调用者无须关心谁来支持调用，系统会自动处理是否调用普通函数或类型的方法。

本节中，首先将用简单的例子了解 Go语言是如何将方法与函数视为一个概念，接着会实现一个事件系统，事件系统能有效地将事件触发与响应两端代码解耦。
方法和函数的统一调用
本节的例子将让一个结构体的方法（class.Do）的参数和一个普通函数（funcDo）的参数完全一致，也就是方法与函数的签名一致。然后使用与它们签名一致的函数变量（delegate）分别赋值方法与函数，接着调用它们，观察实际效果。

详细实现请参考下面的代码。

```go
package main

import "fmt"

// 声明一个结构体
type class struct {
}

// 给结构体添加Do方法
func (c *class) Do(v int) {
	fmt.Println("call method do:", v)
}

// 普通函数的 Do
func FuncDo(v int) {
	fmt.Println("call function do:", v)
}

func main() {
	// 声明一个函数回调
	var delegate func(int)

	// 创建结构体实例返回结构体指针
	c := new(class)

	// 将回调设为c的 Do方法
	delegate = c.Do

	// 调用
	delegate(100)

	// 将回调设为普通函数
	delegate = FuncDo

	// 调用
	delegate(100)
}
```

代码说明如下：

- 第 10 行，为结构体添加一个 Do() 方法，参数为整型。这个方法的功能是打印提示和输入的参数值。
- 第 15 行，声明一个普通函数，参数也是整型，功能是打印提示和输入的参数值。
- 第 21 行，声明一个 delegate 的变量，类型为 func(int)，与 funcDo 和 class 的 Do() 方法的参数一致。
- 第 27 行，将 c.Do 作为值赋给 delegate 变量。
- 第 30 行，调用 delegate() 函数，传入 100 的参数。此时会调用 c 实例的 Do() 方法。
- 第 33 行，将 funcDo 赋值给 delegate。
- 第 36 行，调用 delegate()，传入 100 的参数。此时会调用 funcDo() 方法。

运行代码，输出如下：

```text
call method do: 100
call function do: 100
```

这段代码能运行的基础在于：无论是普通函数还是结构体的方法，只要它们的签名一致，与它们签名一致的函数变量就可以保存普通函数或是结构体方法。

了解了 Go语言的这一特性后，我们就可以将这个特性用在事件中。

#### 事件系统基本原理

事件系统可以将事件派发者与事件处理者解耦。例如，网络底层可以生成各种事件，在网络连接上后，网络底层只需将事件派发出去，而不需要关心到底哪些代码来响应连接上的逻辑。或者再比如，你注册、关注或者订阅某“大V”的社交消息后，“大V”发生的任何事件都会通知你，但他并不用了解粉丝们是如何为她喝彩或者疯狂的。如下图所示为事件系统基本原理图。

![事件系统基本原理](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191104_6.jpg)

一个事件系统拥有如下特性：
- 能够实现事件的一方，可以根据事件 ID 或名字注册对应的事件。
- 事件发起者，会根据注册信息通知这些注册者。
- 一个事件可以有多个实现方响应。

通过下面的步骤详细了解事件系统的构成及使用。

#### 事件注册

事件系统需要为外部提供一个注册入口。这个注册入口传入注册的事件名称和对应事件名称的响应函数，事件注册的过程就是将事件名称和响应函数关联并保存起来，详细实现请参考下面代码的 RegisterEvent() 函数。

```go
package main

import "fmt"

// 实例化一个通过字符串映射函数切片的map
var eventByName = make(map[string][]func(interface{}))

// 注册事件，提供事件名和回调函数
func RegisterEvent(name string, callback func(interface{})) {
	// 通过名字查找事件列表
	list := eventByName[name]

	// 在列表切片中添加函数
	list = append(list, callback)

	// 将修改的事件列表切片保存回去
	eventByName[name] = list
}

// 调用事件
func CallEvent(name string, param interface{}) {
	// 通过名字找到事件列表
	list := eventByName[name]

	// 遍历这个事件的所有回调
	for _, callback := range list {
		// 传入参数调用回调
		callback(param)
	}
}

// 声明角色的结构体
type Actor struct {
}

// 为角色添加一个事件处理函数
func (a *Actor) OnEvent(param interface{}) {
	fmt.Println("actor event:", param)
}

// 全局事件
func GlobalEvent(param interface{}) {
	fmt.Println("global event:", param)
}

func main() {
	// 实例化一个角色
	a := new(Actor)

	// 注册名为 OnSkill 的回调
	RegisterEvent("OnSkill", a.OnEvent)

	// 再次在 OnSkill 上注册全局事件
	RegisterEvent("OnSkill", GlobalEvent)

	CallEvent("OnSkill", 100)
}
```

代码说明如下：
- 第 6 行，创建一个 map 实例，这个 map 通过事件名（string）关联回调列表（[]func(interface{}），同一个事件名称可能存在多个事件回调，因此使用回调列表保存。回调的函数声明为 func(interface{})。
- 第 9 行，提供给外部的通过事件名注册响应函数的入口。
- 第 11 行，eventByName 通过事件名（name）进行查询，返回回调列表（[]func(interface{}）。
- 第 14 行，为同一个事件名称在已经注册的事件回调的列表中再添加一个回调函数。
- 第 17 行，将修改后的函数列表设置到 map 的对应事件名中。

拥有事件名和事件回调函数列表的关联关系后，就需要开始准备事件调用的入口了。

#### 事件调用

事件调用方和注册方是事件处理中完全不同的两个角色。事件调用方是事发现场，负责将事件和事件发生的参数通过事件系统派发出去，而不关心事件到底由谁处理；事件注册方通过事件系统注册应该响应哪些事件及如何使用回调函数处理这些事件。事件调用的详细实现请参考上面代码的 CallEvent() 函数。

代码说明如下：
- 第 21 行，调用事件的入口，提供事件名称 name 和参数 param。事件的参数表示描述事件具体的细节，例如门打开的事件触发时，参数可以传入谁进来了。
- 第 23 行，通过注册事件回调的 eventByName 和事件名字查询处理函数列表 list。
- 第 26 行，遍历这个事件列表，如果没有找到对应的事件，list 将是一个空切片。
- 第 28 行，将每个函数回调传入事件参数并调用，就会触发事件实现方的逻辑处理。

#### 使用事件系统

例子中，在 main() 函数中调用事件系统的 CallEvent 生成 OnSkill 事件，这个事件有两个处理函数，一个是角色的 OnEvent() 方法，还有一个是函数 GlobalEvent()，详细代码实现过程请参考下面的代码。

```go
// 声明角色的结构体
type Actor struct {
}

// 为角色添加一个事件处理函数
func (a *Actor) OnEvent(param interface{}) {
	fmt.Println("actor event:", param)
}

// 全局事件
func GlobalEvent(param interface{}) {
	fmt.Println("global event:", param)
}

func main() {
	// 实例化一个角色
	a := new(Actor)

	// 注册名为 OnSkill 的回调
	RegisterEvent("OnSkill", a.OnEvent)

	// 再次在 OnSkill 上注册全局事件
	RegisterEvent("OnSkill", GlobalEvent)

	CallEvent("OnSkill", 100)
}
```

代码说明如下：

- 第 2 行，声明一个角色的结构体。在游戏中，角色是常见的对象，本例中，角色也是 OnSkill 事件的响应处理方。
- 第 6 行，为角色结构添加一个 OnEvent() 方法，这个方法拥有 param 参数，类型为 interface{}，与事件系统的函数（func(interface{})）签名一致。
- 第 11 行为全局事件响应函数。有时需要全局进行侦听或者处理一些事件，这里使用普通函数实现全局事件的处理。
- 第 20 行，注册一个 OnSkill 事件，实现代码由 a 的 OnEvent 进行处理。也就是 Actor的OnEvent() 方法。
- 第 23 行，注册一个 OnSkill 事件，实现代码由 GlobalEvent 进行处理，虽然注册的是同一个名字的事件，但前面注册的事件不会被覆盖，而是被添加到事件系统中，关联 OnSkill 事件的函数列表中。
- 第 25 行，模拟处理事件，通过 CallEvent() 函数传入两个参数，第一个为事件名，第二个为处理函数的参数。

整个例子运行结果如下：

```go
actor event: 100
global event: 100
```

结果演示，角色和全局的事件会按注册顺序顺序地触发。

一般来说，事件系统不保证同一个事件实现方多个函数列表中的调用顺序，事件系统认为所有实现函数都是平等的。也就是说，无论例子中的 a.OnEvent 先注册，还是 GlobalEvent() 函数先注册，最终谁先被调用，都是无所谓的，开发者不应该去关注和要求保证调用的顺序。

一个完善的事件系统还会提供移除单个和所有事件的方法。





