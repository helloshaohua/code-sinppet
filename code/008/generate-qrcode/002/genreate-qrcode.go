package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func main() {
	engine := gin.New()
	engine.GET("/code", func(context *gin.Context) {
		bytes, err := qrcode.Encode("https://jd.com", qrcode.Highest, 400)
		if err != nil {
			context.Writer.Write([]byte("Generate QRCODE Failed"))
		}
		context.Writer.Write(bytes)
	})
	log.Fatal(engine.Run(":8295"))
}
