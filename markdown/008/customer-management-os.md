### 客户信息管理系统

本节带领大家实现一个基于文本界面的客户关系管理软件，该软件可以实现对客户的插入、修改和删除，并且可以打印客户信息明细表。

软件由一下三个模块组成：

![客户信息管理系统](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191108_6.gif)

项目结构如下所示：

![项目结构](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191108_6.png)

#### 在 costumer.go 中

代码如下：

```go
package model

import "fmt"

// 声明一个 Customer 结构体，表示一个客户信息
type Customer struct {
	ID     int
	Name   string
	Gender string
	Age    int
	Phone  string
	Email  string
}

// 返回客户信息，格式化字符串
func (c Customer) GetInfo() string {
	return fmt.Sprintf(
		"%v\t %v\t %v\t %v\t %v\t %v\t",
		c.ID,
		c.Name,
		c.Gender,
		c.Age,
		c.Phone,
		c.Email,
	)
}

// 使用工厂模式，返回一个 Customer 的实例
func NewCustomer(id int, name, gender string, age int, phone, email string) *Customer {
	return &Customer{
		ID:     id,
		Name:   name,
		Gender: gender,
		Age:    age,
		Phone:  phone,
		Email:  email,
	}
}

// 第二种创建 Customer 实例方法，不带ID
func NewCustomer2(name, gender string, age int, phone, email string) *Customer {
	return &Customer{
		Name:   name,
		Gender: gender,
		Age:    age,
		Phone:  phone,
		Email:  email,
	}
}
```

#### 在 costumerService.go 中

代码如下：

```go
package service

import (
	"code-snippet/code/008/customer-management-os/model"
)

// 该CustomerService，完成对 Customer 的操作, 包括增删改查
type CustomerService struct {
	customers      []*model.Customer
	customerNumber int // 声明一个字段，表示当前切片含有多少个客户该字段后面，还可以作为新客户的 ID+1
}

// 返回客户切片列表
func (c *CustomerService) List() []*model.Customer {
	return c.customers
}

// 添加新客户
func (c *CustomerService) Add(customer *model.Customer) bool {
	// 我们确定一个分配ID 的规则，就是添加的顺序
	c.customerNumber++
	customer.ID = c.customerNumber
	c.customers = append(c.customers, customer)
	return true
}

// 修改客户信息
func (c *CustomerService) Change(id int, name, gender string, age int, phone, email string) bool {
	index := c.FindByID(id)
	if index != -1 {
		customer := c.customers[index]
		if name != "" {
			customer.Name = name
		}
		if gender != "" {
			customer.Gender = gender
		}
		if age != -1 {
			customer.Age = age
		}
		if phone != "" {
			customer.Phone = phone
		}
		if email != "" {
			customer.Email = email
		}
		return true
	}
	return false
}

// 根据 ID 删除客户(从切片中删除)
func (c *CustomerService) Delete(id int) bool {
	index := c.FindByID(id)
	if index != -1 {
		c.customers = append(c.customers[:index], c.customers[index+1:]...)
		return true
	}
	return false
}

// 根据 ID 查找客户在切片中的对应索引值下标，如果没有该客户，返回-1
func (c *CustomerService) FindByID(id int) int {
	index := -1
	// c.customer切片
	for i, customer := range c.customers {
		if id == customer.ID {
			// 找到了
			index = i
		}
	}
	return index
}

// 使用工厂模式构造函数，返回一个 CustomerService 的实例
func NewCustomerService() *CustomerService {
	service := new(CustomerService)
	// 为了能够看到有客户在切片中，初始化时添加一个默认的客户
	service.customerNumber = 1
	service.customers = append(service.customers,
		model.NewCustomer(service.customerNumber, "张三", "男", 20, "010-52328282", "zhangsan@foxmail.com"),
	)
	return service
}
```

#### 在 costumerView.go 中

代码如下：

```go
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
```

#### 在 main.go 中

代码如下：

```go
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
```

#### 执行结果如下所示

```perl
$ go run main.go
----------------------客户信息管理系统-----------------------
1 添加客户
2 修改客户
3 删除客户
4 客户列表
5 退出
请选择(1-5)：1
--------------------------添加客户--------------------------
姓名：武沫汐
性别：男
年龄：22
电话：13729283282
邮箱：wu.moxi@foxmail.com
添加客户成功！
----------------------客户信息管理系统-----------------------
1 添加客户
2 修改客户
3 删除客户
4 客户列表
5 退出
请选择(1-5)：2
--------------------------客户列表--------------------------
编码	姓名	性别	年龄	电话	邮箱	
1	 张三	 男	 20	 010-52328282	 zhangsan@foxmail.com	
2	 武沫汐	 男	 22	 13729283282	 wu.moxi@foxmail.com	
--------------------------修改客户--------------------------
输入客户ID：2
姓名：武沫汐CHANGE
性别：
年龄：
电话：
邮箱：
修改客户成功！
----------------------客户信息管理系统-----------------------
1 添加客户
2 修改客户
3 删除客户
4 客户列表
5 退出
请选择(1-5)：3
--------------------------客户列表--------------------------
编码	姓名	性别	年龄	电话	邮箱	
1	 张三	 男	 20	 010-52328282	 zhangsan@foxmail.com	
2	 武沫汐CHANGE	 男	 22	 13729283282	 wu.moxi@foxmail.com	
--------------------------删除客户--------------------------
请选择待删除的客户编号【-1退出】：1
确认是否删除(Y/N)：Y
删除客户成功！
----------------------客户信息管理系统-----------------------
1 添加客户
2 修改客户
3 删除客户
4 客户列表
5 退出
请选择(1-5)：4
--------------------------客户列表--------------------------
编码	姓名	性别	年龄	电话	邮箱	
2	 武沫汐CHANGE	 男	 22	 13729283282	 wu.moxi@foxmail.com	
----------------------客户信息管理系统-----------------------
1 添加客户
2 修改客户
3 删除客户
4 客户列表
5 退出
请选择(1-5)：5
确认是否退出(Y/N)：Y
已退出客户关系管理系统!
```

#### 项目代码

[customer-management-os](https://github.com/wumoxi/code-sinppet/tree/master/code/008/customer-management-os)