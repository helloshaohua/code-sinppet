package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "hello world")
	})
	log.Fatal(http.ListenAndServe(":8295", nil))
}
