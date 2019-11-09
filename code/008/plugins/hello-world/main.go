package main

import "github.com/dullgiulio/pingo"

// 创建要导出的对象
type MyPlugin struct {
}

// 导出的方法，带有RPC 签名
func (p *MyPlugin) SayHello(name string, msg *string) error {
	*msg = "hello, " + name
	return nil
}

func main() {
	plugin := &MyPlugin{}
	pingo.Register(plugin)
	pingo.Run()
}
