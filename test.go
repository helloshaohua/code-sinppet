package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// 设置 Cookie
	http.HandleFunc("/set-cookie", func(writer http.ResponseWriter, request *http.Request) {
		cookie := &http.Cookie{
			Name:     "username",
			Value:    base64.StdEncoding.EncodeToString([]byte("wumoxi")),
			Expires:  time.Now().AddDate(0, 1, 0),
			HttpOnly: true,
		}
		http.SetCookie(writer, cookie)
		io.WriteString(writer, "hello world")
		return
	})

	// 获取 Cookie
	http.HandleFunc("/get-cookie", func(writer http.ResponseWriter, request *http.Request) {
		// 获取 Cookie
		cookie, err := request.Cookie("username")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		// 解析Cookie值
		bytes, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// 写入响应对象
		_, err = writer.Write(bytes)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		return
	})

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
