package cls2

import (
	"code-snippet/code/008/clsfactory/base"
	"fmt"
)

// 定义类2
type Class2 struct {
}

// 实现Class接口
func (c *Class2) Do() {
	fmt.Println("Class2")
}

func init() {
	// 在启动时注册类2工厂
	base.Register("Class2", func() base.Class {
		return new(Class2)
	})
}
