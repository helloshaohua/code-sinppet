package main

import (
	"fmt"
	"math"
)

type Vec2 struct {
	X float32
	Y float32
}

// 加
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

// 减
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

// 乘
func (v Vec2) Scale(s float32) Vec2 {
	return Vec2{X: v.X * s, Y: v.Y * s}
}

// 距离
func (v Vec2) DistanceTo(other Vec2) float32 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// 差值
func (v Vec2) Normalize() Vec2 {
	mag := v.X*v.X + v.Y*v.Y
	if mag > 0 {
		oneOverMag := 1 / float32(math.Sqrt(float64(mag)))
		return Vec2{X: v.X * oneOverMag, Y: v.Y * oneOverMag}
	}
	return Vec2{0, 0}
}

type Player struct {
	currPos   Vec2    // 当前位置
	targetPos Vec2    // 目标位置
	spend     float32 // 移动速度
}

// 移动到某个点就是设置目标位置
func (p *Player) MoveTo(v Vec2) {
	p.targetPos = v
}

// 获取当前的位置
func (p *Player) Pos() Vec2 {
	return p.currPos
}

// 是否到达
func (p *Player) IsArrived() bool {
	return p.currPos.DistanceTo(p.targetPos) < p.spend
}

// 逻辑更新
func (p *Player) Update() {
	if !p.IsArrived() {
		// 计算出当前位置指向目标的朝向
		dir := p.targetPos.Sub(p.currPos).Normalize()

		// 添加速度矢量生成新的位置
		newPos := p.currPos.Add(dir.Scale(p.spend))

		p.currPos = newPos
	}
}

// 创建新玩家
func NewPlayer(spend float32) *Player {
	return &Player{spend: spend}
}

func main() {
	p := NewPlayer(0.5)
	p.MoveTo(Vec2{3, 1})
	for !p.IsArrived() {
		p.Update()
		fmt.Println(p.Pos())
	}
}
