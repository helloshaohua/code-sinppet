package main

import "fmt"

// 字段结构
type Dictionary struct {
	data map[interface{}]interface{}
}

// 设置键值
func (d *Dictionary) Set(key, value interface{}) {
	d.data[key] = value
}

// 根据键获取值
func (d *Dictionary) Get(key interface{}) interface{} {
	return d.data[key]
}

// 遍历所有的键值，如果回调函数返回值为false，停止遍历
func (d *Dictionary) Visit(callback func(key, value interface{}) bool) {
	if callback == nil {
		return
	}

	for key, value := range d.data {
		if !callback(key, value) {
			return
		}
	}
}

// 清空所有字典数据
func (d *Dictionary) Clear() {
	d.data = make(map[interface{}]interface{})
}

// 创建一个字典结构
func NewDictionary() *Dictionary {
	dictionary := new(Dictionary)
	dictionary.Clear()
	return dictionary
}

func main() {
	// 创建字典实例
	dictionary := NewDictionary()

	// 添加字典数据
	dictionary.Set("My Factory", 60)
	dictionary.Set("Terra Craft", 36)
	dictionary.Set("Don't Hungry", 24)

	// 获取值及打印
	favorite := dictionary.Get("Terra Craft")
	fmt.Println("favorite:", favorite)

	// 遍历所有的字典元素
	dictionary.Visit(func(key, value interface{}) bool {
		// 将值转为int 类型，并判断是否大于40
		if val, ok := value.(int); val > 40 && ok {
			// 输出很贵
			fmt.Println(key, "is expensive")
			return true
		}

		// 默认都是输出很便宜
		fmt.Println(key, "is cheap")
		return true
	})
}
