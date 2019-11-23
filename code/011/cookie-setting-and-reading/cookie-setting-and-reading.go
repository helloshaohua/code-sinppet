package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 获取所有 Cookie
	http.HandleFunc("/get-cookie-all", func(writer http.ResponseWriter, request *http.Request) {
		// 获取所有的 Cookie
		cookies := request.Cookies()
		for _, cookie := range cookies {
			bytes, _ := base64.StdEncoding.DecodeString(cookie.Value)
			fmt.Fprintf(writer, "cookie name: %s, cookie value: %s\n", cookie.Name, string(bytes))
		}
		return
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
