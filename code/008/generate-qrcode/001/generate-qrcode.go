package main

import (
	"log"

	"github.com/skip2/go-qrcode"
)

func main() {
	err := qrcode.WriteFile("https://baidu.com", qrcode.Medium, 256, "./baidu.com.png")
	if err != nil {
		log.Fatal(err)
	}
}
