### Go语言发送电子邮件

电子邮件在日常工作中有很大用途，凡项目或任务，有邮件来往可避免扯皮背锅。而在一些自动化的应用场合，也使用得广泛，特别是系统监控方面，如果在资源使用达到警戒线之前自动发邮件通知运维人员，能消除隐患于前期，而不至于临时临急去做善后方案。

对于多人协合（不管是不是异地）场合，邮件也有用武之地，当有代码或文档更新时，自动发邮件通知项目成员或领导，提醒各方人员知晓并及时更新。

说到发邮件，不得不提用程序的方式实现。下面就来为大家介绍一下怎么使用Go语言来实现发送电子邮件。Go语言拥有大量的库，非常方便使用。

Go语言使用 gomail 包来发送邮箱，代码如下所示：

```go
package main

import (
	"log"
	"toolkit"
)

func main() {
	mailer := toolkit.NewEmail(&toolkit.MailerParams{
		ServerHost:   "smtp.qq.com",
		ServerPort:   465,
		FromEmail:    "shaohua@foxmail.com",
		FromPassword: "mmooqxhsssdothsdpsubddf",
		FromName:     "武沫汐",
		Toers:        []string{"warner@126.com", "warner@139.com", "contact.warner@gmail.com"},
		CCers:        []string{"warner@163.com"},
	})

	send, err := mailer.Send("Golang邮件发送", `中华人民共和国`, "text/plain")
	if err != nil && !send {
		log.Fatal(err)
	} else {
		log.Println("发送成功")
	}
}
```

代码执行结果如下：

```text
2019/11/09 11:56:13 发送成功
```

> 注意: 将具体的邮件操作封装到了工具包中 [toolkit](https://github.com/wumoxi/toolkit)

#### 使用自定义客户端发送邮件需要以下两个要素

1、 发送方的邮箱必须开启 stmt 和 pop3 通道，以 qq 邮箱为例，登陆 qq 邮箱 -> 设置 -> 账户 -> 开启 pop3 和 stmt 服务

![开启 stmt 和 pop3 通道](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191109_8.png)

2、开启后会获得该账户的授权码，如果忘记也可以重新生成。