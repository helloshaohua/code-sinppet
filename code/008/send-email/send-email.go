package main

import (
	"log"
	"toolkit"
)

func main() {
	mailer := toolkit.NewEmail(&toolkit.MailerParams{
		ServerHost:   "smtp.qq.com",
		ServerPort:   465,
		FromEmail:    "wu.shaohua@foxmail.com",
		FromPassword: "mmooqxhothpubddf",
		FromName:     "武沫汐",
		Toers:        []string{"warnerwu@126.com", "warnerwu@139.com", "contact.shaohua@gmail.com"},
		CCers:        []string{"warnerwu@163.com"},
	})

	send, err := mailer.Send("Golang邮件发送", `中华人民共和国`, "text/plain")
	if err != nil && !send {
		log.Fatal(err)
	} else {
		log.Println("发送成功")
	}
}
