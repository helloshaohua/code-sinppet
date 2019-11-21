package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	response, err := http.Get("http://www.baidu.com")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	io.Copy(os.Stdout, response.Body)
}
