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
