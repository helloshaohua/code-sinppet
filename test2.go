package main

import (
	"fmt"
	"sync"
)

type config struct {
}

var (
	cfg  *config
	Once sync.Once
)

func GetConfigInstance() *config {
	if cfg == nil {
		fmt.Println("不存在实例则创建实例对象")
		cfg = new(config)
		return cfg
	}
	fmt.Println("存在实例对象则直接返回")
	return cfg
}

func GetConfigOnceInstance() *config {
	Once.Do(func() {
		fmt.Println("不存在实例则创建实例对象")
		cfg = new(config)
	})
	fmt.Println("存在实例对象则直接返回")
	return cfg
}

func init() {
	cfg = new(config)
}

func NewConfig() *config {
	return cfg
}

func main() {
	fmt.Printf("%T\n", NewConfig())
	fmt.Printf("%T\n", NewConfig())
	fmt.Printf("%T\n", NewConfig())
	fmt.Printf("%T\n", NewConfig())
	fmt.Printf("%T\n", NewConfig())
	fmt.Printf("%T\n", NewConfig())
	fmt.Printf("%T\n", NewConfig())
}
