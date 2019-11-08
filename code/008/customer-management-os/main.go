package main

import (
	"code-snippet/code/008/customer-management-os/service"
	"code-snippet/code/008/customer-management-os/view"
)

func main() {
	// 在 main 函数中，创建一个 customerView，并运行显示主菜单
	customerView := view.CustomerView{
		Key:  "",
		Loop: true,
	}
	// 这里完成对 customerView结构体customerService字段的初始化
	customerView.CustomerService = service.NewCustomerService()
	// 显示主菜单
	customerView.MainMenu()
}
