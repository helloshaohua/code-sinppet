## Golang 代码片段

##### [Go语言输出正弦函数（Sin）图像](markdown/002/out-sin-func-picture.md)

### Go语言字符串

##### [Go语言遍历字符串——获取每一个字符串元素](markdown/002/traversing-string.md)
##### [Go语言字符串截取（获取字符串的某一段字符）](markdown/002/string-interception.md)
##### [Go语言修改字符串](markdown/002/change-string.md)
##### [Go语言字符串拼接（连接）](markdown/002/join-string.md)
##### [Go语言fmt.Sprintf（格式化输出）](markdown/002/format-out-string.md)
##### [Go语言Base64编码——电子邮件的基础编码格式](markdown/002/base64-string.md)

### Go语言指针

##### [使用指针变量获取命令行的输入信息](markdown/002/pointer.md)

### 常量

##### [将枚举值转换为字符串](markdown/002/enumerate-to-string.md)

### 类型别名

##### [在结构体成员嵌入时使用别名](markdown/002/type-alias.md)

### 反射

##### [通过反射获取结构体字段名和字段类型](markdown/002/type-alias.md)

### 流程控制

##### [Go语言输出九九乘法表](markdown/004/for-loop-using.md)
##### [示例：聊天机器人](markdown/004/chat-robot.md)
##### [示例：词频统计](markdown/004/word-frequency-statistics.md)
##### [示例：缩进排序](markdown/004/indentation-sort.md)
##### [示例：二分查找算法](markdown/004/binary-find.md)
##### [示例：冒泡排序](markdown/004/bubble-sort.md)


### Go语言函数

##### [字符串的链式处理](markdown/005/string-chain-processing.md)

### Go语言语法糖

- 函数可变参数列表`func(users ...User)`
- 数组可变长度值 `[...]int`
- 访问结构体指针的成员变量 `player := new(Player) player.Name => (*player).Name`

### Go语言结构体

##### [创建一个HTTP请求](markdown/006/new-http-request.md)
##### [示例：使用事件系统实现事件的响应和处理](markdown/006/event.md)
##### [示例：使用匿名结构体解析JSON数据](markdown/006/anonymous-struct-parse-json-data.md)

### Go语言接口

##### [示例：使用空接口实现可以保存任意值的字典](markdown/007/save-interface-value-into-dictionary.md)

### Go语言包

##### [示例：Go语言工厂模式自动注册](markdown/008/auto-register-factory.md)
##### [示例：Go语言单例模式](markdown/008/single-mode.md)
##### [示例：使用Go语言生成二维码](markdown/008/generate-qrcode.md)
##### [示例：客户信息管理系统](markdown/008/customer-management-os.md)
##### [示例：发送电子邮件](markdown/008/send-email.md)
##### [示例：Pingo插件化开发](markdown/008/pingo-plugin.md)

### Go语言并发

##### [示例：Go语言无缓冲的通道模拟网球比赛](markdown/009/none-cache-channel.md)
##### [示例：Go语言无缓冲的通道模拟接力赛](markdown/009/none-cache-channel-relay.md)
##### [示例：Go语言模拟远程过程调用](markdown/009/mock-rpc.md)
##### [示例：Go语言使用通道响应计时器的事件](markdown/009/use-channel-response-timer-event.md)
##### [示例：Go语言Telnet回音服务器(TCP服务器的基本结构)](markdown/009/telnet.md)
##### [示例：Go语言竞态检测(检测代码在并发环境下可能出现的问题{通过使用原子访问解决})](markdown/009/race-check.md)
##### [示例：Go语言互斥锁(sync.Mutex)和读写互斥锁(sync.RWMutex)](markdown/009/lock.md)
##### [示例：Go语言等待组](markdown/009/wait-group.md)

### Go语言反射

##### [示例：Go语言使用反射修改基本类型值](markdown/010/reflect-operate-basic-type.md)
##### [示例：对结构的反射操作](markdown/010/reflect-get-struct-members-value.md)
##### [示例：通过反射获取指针指向的元素类型](markdown/010/gets-the-type-of-element-the-pointer-to-by-reflection.md)
##### [示例：Go语言通过反射获取结构体的成员类型](markdown/010/gets-the-member-type-of-the-structure-through-reflection.md)
##### [示例：Go语言结构体标签](markdown/010/struct-field-tag.md)
##### [示例：Go语言通过反射获取值信息](markdown/010/get-value-information-by-reflection.md)
##### [示例：Go语言使用reflect.Type显示一个类型的方法集](markdown/010/get-type-methods-set.md)
##### [示例：Go语言通过反射访问结构体成员的值](markdown/010/values-of-structure-members-are-accessed-by-reflection.md)
##### [示例：判断反射值的空和有效性](markdown/010/determines-the-null-and-validity-of-the-reflection-value.md)
##### [示例：通过反射修改变量的值](markdown/010/modify-the-value-of-a-variable-by-reflection.md)
##### [示例：通过类型信息创建实例](markdown/010/create-an-instance-with-type-information.md)
##### [示例：通过反射调用函数](markdown/010/functions-are-called-by-reflection.md)
##### [示例：依赖注入](markdown/010/injection.md)

### Go语言网络编程

##### [示例：Dial()函数](markdown/011/dial-func.md)
##### [示例：建立TCP链接](markdown/011/setting-up-tcp-links.md)
##### [示例：Go语言DialTCP()](markdown/011/dial-tcp.md)
##### [示例：RPC协议远程过程调用](markdown/011/rpc-protocol.md)
##### [示例：如何设计优雅的RPC接口](markdown/011/how-to-design-an-elegant-rpc-interface.md)
##### [示例：解码未知结构的JSON数据](markdown/011/decoding-json-data-with-an-unknown-structure.md)
##### [示例：开发一个简单的相册网站](markdown/011/photos.md)
##### [示例：并发时钟服务器](markdown/011/clock-server.md)
##### [示例：Cookie的设置与读取](markdown/011/cookie-setting-and-rending.md)
##### [示例：获取IP地址和域名解析](markdown/011/get-ip-address-and-domain-name-resolution.md)
