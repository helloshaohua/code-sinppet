package main

import "fmt"

type Point struct {
	X int
	Y int
}

type Player struct {
	Name        string
	HealthPoint int
	MagicPoint  int
}

func main() {
	p := new(Point)
	p.X = 33
	p.Y = 55

	player := new(Player)
	player.Name = "张三"
	player.HealthPoint = 100
	player.MagicPoint = 2000

	fmt.Printf("%+v", player)
}
