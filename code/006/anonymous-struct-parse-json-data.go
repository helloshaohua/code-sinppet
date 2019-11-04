package main

import (
	"encoding/json"
	"fmt"
)

// 定义手机屏幕
type Screen struct {
	Size       float32
	ResX, ResY int
}

// 定义电池
type Battery struct {
	Capacity int // 容量
}

// 生成 JSON 数据
func generateJSONData() []byte {
	// 完整的数据结构
	raw := &struct {
		Screen
		Battery
		HasTouchID bool // 序列化时添加的字段：是否有指纹识别
	}{
		// 屏幕参数
		Screen: Screen{
			Size: 5.5,
			ResX: 1920,
			ResY: 1080,
		},
		Battery: Battery{
			Capacity: 10000,
		},
		HasTouchID: true,
	}

	bytes, _ := json.Marshal(raw)

	return bytes
}

func main() {
	// 生成一段JSON 数据
	data := generateJSONData()

	// 转换为字符串格式并打印输出
	fmt.Println(string(data))

	// 只需要屏幕和指纹识别信息的结构和实例
	screenAndTouch := struct {
		Screen
		HasTouchID bool
	}{}

	// 反序列化到screenAndTouch中
	json.Unmarshal(data, &screenAndTouch)
	// 输出screenAndTouch的详细结构
	fmt.Printf("screenAndTouch: %+v\n", screenAndTouch)

	// 只需要电池和指纹识别信息的结构和实例
	batteryAndTouch := struct {
		Battery
		HasTouchID bool
	}{}

	// 反序列化到batteryAndTouch
	json.Unmarshal(data, &batteryAndTouch)
	// 输出screenAndTouch的详细结构
	fmt.Printf("batteryAndTouch: %+v\n", batteryAndTouch)
}
