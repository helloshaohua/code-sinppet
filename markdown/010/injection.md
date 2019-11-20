### Go语言inject库：依赖注入

前面己经对反射的基本概念和相关 API 进行了讲解，本节结合一个非常著名的包 inject 展开讲解，inject 借助反射提供了对 2 种类型实体的注入：函数和结构。Go 著名的 Web 框架 martini 的依赖注入使用的就是这个包。

#### 依赖注入和控制反转

在介绍 inject 之前先简单介绍“依赖注入”和“控制反转”的概念。正常情况下，对函数或方法的调用是调用方的主动直接行为，调用方清楚地知道被调的函数名是什么，参数有哪些类型，直接主动地调用；包括对象的初始化也是显式地直接初始化。

所谓的“控制反转”就是将这种主动行为变成间接的行为，主调方不是直接调用函数或对象，而是借助框架代码进行间接的调用和初始化，这种行为我们称为“控制反转”，控制反转可以解藕调用方和被调方。

“库”和“框架”能很好地解释“控制反转”的概念。一般情况下，使用库的程序是程序主动地调用库的功能，但使用框架的程序常常由框架驱动整个程序，在框架下写的业务代码是被框架驱动的，这种模式就是“控制反转”。

“依赖注入”是实现“控制反转”的一种方法，如果说“控制反转”是一种设计思想，那么“依赖注入”就是这种思想的一种实现，通过注入的参数或实例的方式实现控制反转。如果没有特殊说明，我们通常说的“依赖注入”和“控制反转”是一个东西。

大家可能会疑惑，为什么不直接光明正大地调用，而非要拐弯抹角地进行间接调用，控制反转的价值在哪里呢？一句话“解耦”，有了控制反转就不需要调用者将代码写死，可以让控制反转的框架代码读取配置，动态地构建对象，这一点在 Java 的 Spring 框架中体现得尤为突出。

控制反转是解决复杂问题一种方法，特别是在 Web 框架中为路由和中间件的灵活注入提供了很好的方法。但是软件开发没有银弹，当问题足够复杂时，应该考虑的是服务拆分，而不是把复杂的逻辑用一个“大盒子”装起来，看起来干净了，但也只是看起来干净，实现还是很复杂，这也是使用框架带来的副作用。

#### inject 实践

inject 是 Go语言依赖注入的实现，它实现了对结构（struct）和函数的依赖注入。在介绍具体实现之前，先来想一个问题，如何通过一个字符串类型的函数名调用函数。Go 没有 Java 中的 Class.forName 方法可以通过类名直接构造对象，所以这种方法是行不通的，能想到的方法就是使用 map 实现一个字符串到函数的映射，代码如下：

```go
func fl() {
    println ("fl")
}
func f2 () {
    println ("f2")
}
funcs := make(map[string] func ())
funcs ["fl"] = fl
funcs ["f2"] = fl
funcs ["fl"]()
funcs ["f2"]()
```

但是这有个缺陷，就是 map 的 Value 类型被写成 func()，不同参数和返回值的类型的函数并不能通用。将 map 的 Value 定义为 interface{} 空接口类型是否能解决该问题？可以解决该问题，但需要借助类型断言或反射来实现，通过类型断言实现等于又绕回去了，反射是一种可行的办法。inject 包借助反射实现函数的注入调用，下面通过一个例子来看一下。

#### 函数的注入

```go
package main

import (
	"fmt"

	"github.com/go-macaron/inject"
)

type S1 interface {
}

type S2 interface {
}

func Format(name string, company S1, level S2, age int) {
	fmt.Printf("name: %s, company: %s, level: %d, age: %d\n", name, company, level, age)
}

func main() {
	// 控制实例的创建
	i := inject.New()

	// 实参注入
	i.Map("张三")
	i.MapTo("阿里巴巴", (*S1)(nil))
	i.MapTo(8, (*S2)(nil))
	i.Map(18)

	// 函数反转调用
	i.Invoke(Format)
}
```

执行结果：

```text
name: 张三, company: 阿里巴巴, level: 8, age: 18
```

可见 inject 提供了一种注入参数调用函数的通用功能，inject.New() 相当于创建了一个控制实例，由其来实现对函数的注入调用。inject 包不但提供了对函数的注入，还实现了对 struct 类型的注入，看下一个示例。

#### 结构体类型的注入

```go
package main

import (
	"fmt"

	"github.com/go-macaron/inject"
)

type S1 interface {
}

type S2 interface {
}

type Staff struct {
	Name    string `inject:"name"`
	Company S1     `inject:"company"`
	Level   S2     `inject:"level"`
	Age     int    `inject:"age"`
}

func main() {
	// 创建被注入实例
	s := Staff{}

	// 控制实例的创建
	si := inject.New()

	// 初始化注入值
	si.Map("张三")
	si.MapTo("阿里巴巴", (*S1)(nil))
	si.MapTo(8, (*S2)(nil))
	si.Map(18)

	// 实现对 struct 注入
	si.Apply(&s)

	// 打印注入后结构体
	fmt.Printf("s: %+v\n", s)
}
```
执行结果：

```text
s: {Name:张三 Company:阿里巴巴 Level:8 Age:18}
```

可以看到 inject 提供了一种对结构类型的通用注入方法。至此，我们仅仅从宏观层面了解 iniect 能做什么，下面从源码实现角度来分析 inject。

#### inject 原理分析

inject 包只有 187 行代码（包括注释），却提供了一个完美的依赖注入实现，下面采用自顶向下的方法分析其实现原理。

##### 入口函数 New

inject.New() 函数构建一个具体类型 injector 实例作为内部注入引擎，返回的是一个 Injector 类型的接口。这里也体现了一种面向接口的设计思想：对外暴露接口方法，对外隐藏内部实现。示例如下：

```go
func New() Injector {
    return &injector {
        values : make(map[reflect.Type)reflect.Value),
    }
}
```

#### 接口设计

下面来看一下具体的接口设计，Injector 暴露了所有方法给外部使用者，这些方法又可以归纳为两大类。第一类方法是对参数注入进行初始化，将结构类型的字段的注入和函数的参数注入统一成一套方法实现；第二类是专用注入实现，分别是生成结构对象和调用函数方法。

在代码设计上，inject 将接口的粒度拆分得很细，将多个接口组合为一个大的接口，这也符合 Go 的 Duck 类型接口设计的原则。injector 按照上述方法拆分为三个接口。示例如下：

```go
type Injector interface {
    //抽象生成注入结构实例的接口
    Applicator
    //抽象函数调用的接口
    Invoker
    //抽象注入参数的接口
    TypeMapper
    //实现一个注入实例链， 下游的能覆盖上游的类型
    SetParent(Injector)
}
```

TypeMapper 接口实现对注入参数操作的汇总，包括设置和查找相关的类型和值的方法。注意：无论函数的实参，还是结构的字段，在 inject 内部，都存放在 map[reflect.Type]reflect.Value 类型的 map 里面，具体实现在后面介绍 injector 时会讲解。

```go
type TypeMapper interface {
    //如下三个方法是设直参数
    Map(interface{}) TypeMapper
    MapTo(interface{}, interface{}) TypeMapper
    Set(reflect.Type, reflect.Value) TypeMapper
    //查找参数
    Get(reflect.Type) reflect.Value
}
```

Invoker 接口中 Invoke 方法是对被注入实参函数的调用：

```go
type Invoker interface {
    Invoke (interface{}) ([]reflect.Value, error)
}
```

Applicator 接口中 Apply 方法实现对结构的注入：

```go
type Applicator interface {
    Apply(interface{}) error
}
```

下面梳理了整个 inject 包的处理流程：

- 通过 inject.New() 创建注入引擎，注入引擎被隐藏，返回的是 Injector 接口类型变量。
- 调用 TypeMapper 接口（Injector 内嵌 TypeMapper）的方法注入 struct 的字段值或函数的实参值。
- 调用 Invoker 方法执行被注入的函数，或者调用 Applicator 接口方法获得被注入后的结构实例。

##### 内部实现

下面具体看一下 inject 内部注入引擎 injector 的实现，首先看一下 injector 的数据结构。


```go
type injector struct {
    values map[reflect.Type]reflect.Value
    parent Injector
}
```

values 里面存放的可以是被注入 struct 的字段类型和值，也可以是函数实参的类型和值。注意：values 是以 reflect.Type 为 Key 的 map，如果一个结构的字段类型相同，则后面注入的参数会覆盖前面的参数，规避办法是使用 MapTo 方法，通过抽象出一个接口类型来避免被覆盖。

```go
func (i *injector) MapTo (val interface{}, ifacePtr interface{}) TypeMapper {
    i.values[InterfaceOf(ifacePtr)] = reflect.ValueOf (val)
    return i
}
```

injector 里面的 parent 的作用是实现多个注入引擎，其构成了一个链。

下面重点分析 injector 对函数的注入实现。示例如下：

```go
func (inj *injector) Invoke(f interface{}) ([]reflect.Value, error) {
    //获取函数类型的 Type
    t := reflect.TypeOf(f)
    //构造一个存放函数实参 Value 值的数纽
    var in = make([]reflect.Value, t.NumIn())
    //使用反射获取函数实参 reflect.Type，逐个去 injector 中查找注入的 Value 值
    for i := O; i < t.NumIn(); i++ {
        argType := t.In(i)
        val := inj.Get(argType)
        if !val.IsValid() {
            return nil, fmt.Errorf("Value not found for type %v", argType)
        }
        in[i] = val
    }
    //反射调用函数
    return reflect.ValueOf(f).Call(in), nil
}
```

inject 对函数注入调用实现很简洁，就是从 injector 里面获取函数实参，然后调用函数。

通过对 inject 包的分析，认识到其“短小精悍”、功能强大，这些实现的基础是依靠反射。但同时注意到包含反射的代码相对来说复杂难懂，虽然 inject 的实现只有短短 200 行代码，但阅读起来并不是很流畅。所以说反射是一把双刃剑，好用但代码不好读。
