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
