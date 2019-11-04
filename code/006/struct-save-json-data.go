package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	// 声明技能结构
	type Skill struct {
		Name  string
		Level int
	}

	// 声明角色结构
	type Actor struct {
		Name string
		Age  int

		Skills []Skill
	}

	// 填充基本角色数据
	a := Actor{
		Name: "cow boy",
		Age:  37,

		Skills: []Skill{
			Skill{
				Name:  "Roll and roll",
				Level: 1,
			},
			Skill{
				Name:  "Flash your dog eye",
				Level: 2,
			},
			Skill{
				Name:  "Time to have Lunch",
				Level: 3,
			},
		},
	}

	if result, err := MarshalJSON(a); err != nil {
		fmt.Printf("json.Marshal error:%+v\n", err)
	} else {
		fmt.Println(string(result))
	}
}

func MarshalJSON(data interface{}) (string, error) {
	// 准备一个缓冲区
	buffer := new(bytes.Buffer)

	// 将任意值转换为 JSON 并输出到缓冲
	if err := writeAny(buffer, reflect.ValueOf(data)); err != nil {
		return "", err
	} else {
		return buffer.String(), err
	}
}

func writeAny(buffer *bytes.Buffer, value reflect.Value) error {
	switch value.Kind() {
	case reflect.String:
		// 写入带有双引号括起来的字符串
		buffer.WriteString(strconv.Quote(value.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 将整型转换为字符串写入缓冲
		buffer.WriteString(strconv.FormatInt(value.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buffer.WriteString(strconv.FormatInt(value.Int(), 10))
	case reflect.Slice:
		return writeSlice(buffer, value)
	case reflect.Struct:
		return writeStruct(buffer, value)
	case reflect.Map:
		return writeMap(buffer, value)
	default:
		return errors.New("unsupport kind: " + value.Kind().String())
	}
	return nil
}

// 将结构体转换为 JSON 并输出到缓冲区
func writeStruct(buffer *bytes.Buffer, value reflect.Value) error {
	// 获取值的类型对象
	valueType := value.Type()

	// 写入结构体左大括号
	buffer.WriteString("{")

	// 遍历结构体的所有值
	for i := 0; i < value.NumField(); i++ {
		// 获取每个字段的字段值(reflect.value)
		fieldValue := value.Field(i)

		// 获取每个字段的类型(reflect.StructField)
		fieldType := valueType.Field(i)

		// 写入字段名左双引号
		buffer.WriteString("\"")

		// 写入字段名
		buffer.WriteString(fieldType.Name)

		// 写入字段名右双引号和冒号
		buffer.WriteString("\":")

		// 写入每个字段值
		writeAny(buffer, fieldValue)

		// 写入每个字段尾部逗号，最后一个字段不添加
		if i < value.NumField()-1 {
			buffer.WriteString(",")
		}
	}

	// 写入结构体右大括号
	buffer.WriteString("}")

	return nil
}

// 将切片转换为 JSON 并输出到缓冲区
func writeSlice(buffer *bytes.Buffer, value reflect.Value) error {
	// 写入切片开始标记
	buffer.WriteString("[")

	// 遍历每个切片元素
	for i := 0; i < value.Len(); i++ {
		sliceValue := value.Index(i)

		// 写入每个切片元素
		writeAny(buffer, sliceValue)

		// 写入每个元素尾部逗号，最后一个字段不添加
		if i < value.Len()-1 {
			buffer.WriteString(",")
		}
	}

	// 写入切片结束标记
	buffer.WriteString("]")

	return nil
}

// 将映射转换为 JSON 并输出到缓冲区
func writeMap(buffer *bytes.Buffer, value reflect.Value) error {
	// 写入映射开始标记
	buffer.WriteString("{")

	// 获取映射迭代器
	iter := value.MapRange()

	// 获取映射长度
	l := value.Len()

	// 迭代次数计数器
	c := new(int)

	for iter.Next() {
		// 获取每个键
		key := iter.Key().String()

		// 获取每个值
		val := iter.Value()

		// 写入字段名左双引号
		buffer.WriteString("\"")

		// 写入键名
		buffer.WriteString(key)

		// 写入字段名右双引号和冒号
		buffer.WriteString("\":")

		// 写入值
		writeAny(buffer, val)

		// 迭代器累加1，再次判断有没有迭代到最后一个元素，到达最后一个元素时，不添加「，」
		if *c++; *c != l {
			buffer.WriteString(",")
		}
	}
	// 写入映射开始标记
	buffer.WriteString("}")

	return nil
}
