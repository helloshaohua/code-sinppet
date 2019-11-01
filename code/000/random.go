package main

import (
	"math/rand"
	"time"
)

// 生成指定范围的随机数
func Random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
