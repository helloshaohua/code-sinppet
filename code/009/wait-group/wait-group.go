package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	// 声明一个等待组
	var wg sync.WaitGroup

	// 准备一系列的网站地址
	var urls = []string{
		"https://baidu.com",
		"https://taobao.com",
		"https://jd.com",
	}

	// 遍历这些地址
	for _, url := range urls {
		// 每一个任务开始时，将等待组增加1
		wg.Add(1)

		// 开启一个并发
		go func(url string) {
			// 使用defer，表示函数完成时将等待组值减1
			defer wg.Done()

			// 使用http访问提供的地址
			_, err := http.Get(url)

			// 访问完成后，打印地址和可能发生的错误
			fmt.Println(url, err)
		}(url)
	}

	// 等待所有的并发任务完成
	wg.Wait()

	// 输出结束提示信息
	fmt.Println("Get all urls done!")
}
