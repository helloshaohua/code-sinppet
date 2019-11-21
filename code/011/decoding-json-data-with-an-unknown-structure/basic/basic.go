package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	b := []byte(`{
    	"Title": "Go语言编程",
    	"Authors": ["XuShiwei", "HughLv", "Pandaman", "GuaguaSong", "HanTuo", "BertYuan",
    	"XuDaoli"],
    	"Publisher": "ituring.com.cn",
    	"IsPublished": true,
    	"Price": 9.99,
    	"Sales": 1000000
	}`)

	var r interface{}

	err := json.Unmarshal(b, &r)
	if err != nil {
		log.Fatal(err)
	}

	book, ok := r.(map[string]interface{})
	if ok {
		for k, v := range book {
			switch value := v.(type) {
			case string:
				fmt.Printf("【%s】 is string 【%s】\n", k, value)
			case bool:
				fmt.Printf("【%s】 is bool 【%t】\n", k, value)
			case int:
				fmt.Printf("【%s】 is int 【%d】\n", k, value)
			case float64:
				fmt.Printf("【%s】 is float64 【%f】\n", k, value)
			case []interface{}:
				fmt.Printf("%s is an slice:\n", k)
				for i, sv := range value {
					fmt.Printf("index: %d is value of %+v\n", i, sv)
				}
			default:
				fmt.Printf("%s is another type not handle yet.\n", k)
			}
		}
	}
}
