package main

import (
	"image/color"
	"log"

	"github.com/skip2/go-qrcode"
)

func main() {
	code, err := qrcode.New("https://imooc.com", qrcode.Highest)
	if err != nil {
		log.Fatal(err)
	}
	code.BackgroundColor = color.RGBA{50, 205, 50, 255}
	code.ForegroundColor = color.White
	code.WriteFile(400, "./imooc.png")
}
