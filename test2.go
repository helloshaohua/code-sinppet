package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	client := new(http.Client)

	// 创建一个http请求
	request, err := http.NewRequest("POST", "http://www.163.com/", strings.NewReader("key=value"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 为请求头添加信息
	request.Header.Add("User-Agent", "myClient")

	// 开始请求
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 关闭响应体
	defer response.Body.Close()

	// 读取响应体数据
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 输出响应体
	fmt.Println(string(bytes))
}
