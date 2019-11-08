package view

import (
	"code-snippet/code/008/customer-management-os/model"
	"code-snippet/code/008/customer-management-os/service"
	"fmt"
)

type CustomerView struct {
	Key             string                   // 接收用户输入
	Loop            bool                     // 表示是否循环的显示菜单
	CustomerService *service.CustomerService // 增加一个字段 CustomerService
}

// 显示所有的客户信息
func (c *CustomerView) list() {
	// 首先，获取到当前所有的客户信息(在切片中)
	customers := c.CustomerService.List()
	// 显示
	fmt.Println("--------------------------客户列表--------------------------")
	fmt.Println("编码\t姓名\t性别\t年龄\t电话\t邮箱\t")
	for i := 0; i < len(customers); i++ {
		fmt.Println(customers[i].GetInfo())
	}
}

// 获取用户输入信息构建新的客户，并完成添加
func (c *CustomerView) add() {
	fmt.Println("--------------------------添加客户--------------------------")

	fmt.Print("姓名：")
	name := ""
	fmt.Scanln(&name)

	fmt.Print("性别：")
	gender := ""
	fmt.Scanln(&gender)

	fmt.Print("年龄：")
	age := 0
	fmt.Scanln(&age)

	fmt.Print("电话：")
	phone := ""
	fmt.Scanln(&phone)

	fmt.Print("邮箱：")
	email := ""
	fmt.Scanln(&email)

	// 构建一个新的 Customer 实例，
	// 注意：ID号，没有让用户输入，ID 是唯一的，需要系统分配
	customer := model.NewCustomer2(name, gender, age, phone, email)

	// 调用添加方法
	if c.CustomerService.Add(customer) {
		fmt.Println("添加客户成功！")
	} else {
		fmt.Println("添加客户失败!")
	}
}

// 获取用户输入信息构建新的客户，并完成添加
func (c *CustomerView) change() {
	c.list()
	fmt.Println("--------------------------修改客户--------------------------")
	fmt.Print("输入客户ID：")
	id := 0
	fmt.Scanln(&id)

	fmt.Print("姓名：")
	name := ""
	fmt.Scanln(&name)

	fmt.Print("性别：")
	gender := ""
	fmt.Scanln(&gender)

	fmt.Print("年龄：")
	age := -1
	fmt.Scanln(&age)

	fmt.Print("电话：")
	phone := ""
	fmt.Scanln(&phone)

	fmt.Print("邮箱：")
	email := ""
	fmt.Scanln(&email)

	// 调用添加方法
	if c.CustomerService.Change(id, name, gender, age, phone, email) {
		fmt.Println("修改客户成功！")
	} else {
		fmt.Println("修改客户失败!")
	}
}

// 获取用户输入的 ID, 删除该ID对应的客户
func (c *CustomerView) delete() {
	c.list()
	fmt.Println("--------------------------删除客户--------------------------")
confirmation:
	fmt.Print("请选择待删除的客户编号【-1退出】：")
	id := -1
	fmt.Scanln(&id)
	if id == -1 {
		// 放弃删除操作
		return
	}

	for {
		fmt.Println("确认是否删除(Y/N)：")
		choice := ""
		fmt.Scanln(&choice)
		if choice == "Y" || choice == "y" {
			// 调用CustomerService的Delete方法
			if c.CustomerService.Delete(id) {
				fmt.Println("删除客户成功！")
			} else {
				fmt.Println("删除客户失败，输入的ID号不存在")
			}
			break
		}

		if choice == "N" || choice == "n" {
			goto confirmation
		}
	}
}

// 退出软件
func (c *CustomerView) exit() {
	fmt.Print("确认是否退出(Y/N)：")

	for {
		fmt.Scanln(&c.Key)
		if c.Key == "Y" || c.Key == "y" || c.Key == "N" || c.Key == "n" {
			break
		}
		fmt.Print("你的输入有误，确认是否退出(Y/N)：")
	}

	if c.Key == "Y" || c.Key == "y" {
		c.Loop = false
	}
}

// 显示主菜单
func (c *CustomerView) MainMenu() {
	for {
		fmt.Println("----------------------客户信息管理系统-----------------------")
		fmt.Println("1 添加客户")
		fmt.Println("2 修改客户")
		fmt.Println("3 删除客户")
		fmt.Println("4 客户列表")
		fmt.Println("5 退出")
		fmt.Print("请选择(1-5)：")

		fmt.Scanln(&c.Key)
		switch c.Key {
		case "1":
			c.add()
		case "2":
			c.change()
		case "3":
			c.delete()
		case "4":
			c.list()
		case "5":
			c.exit()
		}

		if !c.Loop {
			break
		}
	}
	fmt.Println("已退出客户关系管理系统!")
}
