### Go语言开发一个简单的相册网站

本节我们将综合之前介绍的网站开发相关知识，一步步介绍如何开发一个虽然简单但五脏俱全的相册网站。

#### 新建工程

首先创建一个用于存放工程源代码的目录并切换到该目录中去，随后创建一个名为 photos.go 的文件，用于后面编辑我们的代码：

```perl
$ mkdir -p photos/uploads && cd photos && touch photos.go
```

当前目录结构如下：

```go
photos
├── photos.go
└── uploads
```

我们的示例程序不是再造一个 Flickr 那样的网站或者比其更强大的图片分享网站，虽然我们可能很想这么玩。不过还是先让我们快速开发一个简单的网站小程序，暂且只实现以下最基本的几个功能：

- 支持图片上传；
- 在网页中可以查看已上传的图片；
- 能看到所有上传的图片列表；
- 可以删除指定的图片。

功能不多，也很简单。在大概了解上一节中的网页输出 Hello world 示例后，想必你已经知道可以引入 net/http 包来提供更多的路由分派并编写与之对应的业务逻辑处理方法，只不过会比输出一行 Hello, world! 多一些环节，还有些细节需要关注和处理。

#### 使用 net/http 包提供网络服务

接下来，我们继续使用 Go 标准库中的 net/http 包来一步步构建整个相册程序的网络服务。

##### 1、上传图片

先从最基本的图片上传着手，具体代码如下所示。

```go
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// 图片上传
func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	// 显示表单
	if request.Method == "GET" {
		io.WriteString(writer, `<!DOCTYPE html>
			<html lang="en">
			<head>
			    <meta charset="UTF-8">
			    <title>Title</title>
			</head>
			<body>
			<form method="post" action="/upload" enctype="multipart/form-data">
			    Choose an image to upload:
			    <input name="image" type="file">
			    <input value="Upload" type="submit">
			</form>
			</body>
			</html>
		`)
		return
	}
}
```

可以看到，结合 main() 和 uploadHandler() 方法，针对 HTTP GET 方式请求 /upload 路径，程序将会往 http.ResponseWriter 类型的实例对象 writer 中写入一段 HTML 文本，即输出一个 HTML 上传表单。

如果我们使用浏览器访问这个地址，那么网页上将会是一个可以上传文件的表单。

![上传文件表单](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191121_1.png)

光有上传表单还不能完成图片上传，服务端程序还必须有接收上传图片的相关处理。针对上传表单提交过来的文件，我们对 uploadHandler() 方法再添加些业务逻辑程序：

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	UploadDir = "./code/011/photos/uploads"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// 图片上传
func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	// 显示表单
	if request.Method == "GET" {
		io.WriteString(writer, `<!DOCTYPE html>
			<html lang="en">
			<head>
			    <meta charset="UTF-8">
			    <title>Title</title>
			</head>
			<body>
			<form method="post" action="/upload" enctype="multipart/form-data">
			    Choose an image to upload:
			    <input name="image" type="file">
			    <input value="Upload" type="submit">
			</form>
			</body>
			</html>
		`)

		return
	}

	// 上传图片
	if request.Method == "POST" {
		// 获取上传文件
		file, header, err := request.FormFile("image")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v\n", header)

		// 获取上传文件的文件名
		filename := header.Filename

		// 延迟关闭文件
		defer file.Close()

		// 生成新文件名
		nfn := getNewFileNameForUpload(filename)

		// 创建文件
		create, err := os.Create(getFilePath(UploadDir, nfn))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// 延迟关闭新创建的文件
		defer create.Close()

		// 将上传的文件内容拷贝到新创建的文件
		if _, err = io.Copy(create, file); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		// 重定向到预览地址
		http.Redirect(writer, request, "/view?id="+nfn, http.StatusFound)

		// 返回停止处理
		return
	}
}

// 获取新的上传文件名
func getNewFileNameForUpload(filename string) string {
	fs := strings.Split(filename, ".")
	ext := strings.ToLower(fs[len(fs)-1])
	t := time.Unix(time.Now().Unix(), 0).Format("2006-01-02-15-04-05")
	m := strconv.FormatInt(time.Now().UnixNano(), 10)
	return t + m + "." + ext
}

// 获取文件路径
func getFilePath(dir, filename string) string {
	return dir + "/" + filename
}
```

如果是客户端发起的 HTTP POST 请求，那么首先从表单提交过来的字段中寻找名为 image 的文件域并对其接收，调用 request.FormFile() 方法会返回 3 个值，各个值的类型分别是 multipart.File、*multipart.FileHeader 和 error。

如果上传的图片接收不成功，那么在示例程序中返回一个 HTTP 服务端的内部错误给客户端。如果上传的图片接收成功，则将该图片的内容复制到一个临时文件里。如果临时文件创建失败，或者图片副本保存失败，都将触发服务端内部错误。

如果临时文件创建成功并且图片副本保存成功，即表示图片上传成功，就跳转到查看图片页面。此外，我们还定义了两个 defer 语句，无论图片上传成功还是失败，当 uploadHandler() 方法执行结束时，都会先关闭临时文件句柄，继而关闭图片上传到服务器文件流的句柄。

我们还添加了2个辅助函数 `getNewFileNameForUpload` 、`getFilePath` 分别用来通过已有文件名生成新文件名和获取文件路径。

当图片上传成功后，我们即可在网页上查看这张图片，顺便确认图片是否真正上传到了服务端。接下来在网页中呈现这张图片。

##### 2、在网页上显示图片

要在网页中显示图片，必须有一个可以访问到该图片的网址。在前面的示例代码中，图片上传成功后会跳转到 /view?id=filename 这样的网址，因此我们的程序要能够将对 /view 路径的访问映射到某个具体的业务逻辑处理方法。

首先，在 photos.go 程序中新增一个名为 viewHandler 的方法，其代码如下：

```go
// 预览图片
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	// 目标文件
	dst := getFilePath(UploadDir, request.FormValue("id"))

	// 获取文件类型
	contentType, err := getFileContentType(dst)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置Header头信息
	writer.Header().Set("Content-Type", contentType)

	// 从服务器读取文件并作为响应数据输出给客户端
	http.ServeFile(writer, request, dst)

	// 返回停止处理
	return
}

// 获取文件MIME类型
func getFileContentType(filepath string) (contentType string, err error) {
	// 打开文件
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 只需要读取文件的前512个字节就够了
	buffer := make([]byte, 512)

	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}
```

在上述代码中，我们首先从客户端请求中对参数进行接收。request.FormValue("id") 即可得到客户端请求传递的图片唯一ID(或者说文件名)，然后我们将图片 ID 结合之前保存图片用的目录进行组装，即可得到文件在服务器上的存储路径。

接着，调用 http.ServeFile() 方法将该路径下的文件从磁盘中读取并作为服务端的响应数据输出给客户端。同时，也将 HTTP 响应头输出格式设置为通过 getFileContentType 获取的具体的MIME类型。

完成 viewHandler 的业务逻辑后，我们将该方法注册到程序的 main() 方法中，与 /view 访问路径形成映射关联。main() 方法的代码如下：

```go
func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/view", viewHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

这样当客户端(浏览器)访问 /view 路径并传递 id 参数时，即可直接以 HTTP 形式看到图片的内容。在网页上，将会呈现一张可视化的图片。

接下来可以访问 `localhost:8080/upload`，展示一个上传表单，选择要上传的文件，点击 `Upload` 按钮。即可完成图片的上传，并跳转到 `localhost:8080/view?id=xxxxx.jpg` 类似的页面，也可以成功预览查看上传后的图像。

![上传并预览上传](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191121_2.gif)

##### 3、处理不存在的图片访问

理论上，只要是 uploads/ 目录下有的图片，都能够访问到，但我们还是假设有意外情况，比如网址中传入的图片 ID 在 uploads/ 没有对应的文件，这时，我们的 viewHandler() 方法就显得很脆弱了。

不管是给出友好的错误提示还是返回 404 页面，都应该对这种情况作相应处理。我们不妨先以最简单有效的方式对其进行处理，修改 viewHandler() 方法，具体如下：

```go
// 预览图片
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	// 目标文件
	dst := getFilePath(UploadDir, request.FormValue("id"))

	// 检测文件是否存在，不存在则抛出404
	if exists := isExists(dst); !exists {
		http.NotFound(writer, request)
		return
	}

	// 获取文件类型
	contentType, err := getFileContentType(dst)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置Header头信息
	writer.Header().Set("Content-Type", contentType)

	// 从服务器读取文件并作为响应数据输出给客户端
	http.ServeFile(writer, request, dst)

	// 返回停止处理
	return
}

// 检测文件是否存在
func isExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
```

同时，我们增加了 isExists() 辅助函数，用于检查文件是否真的存在。

##### 4、列出所有已上传的图片

这个程序应该有个入口，可以看到所有已上传的图片。对于所有列出的这些图片，我们可以选择进行查看或者删除等操作。下面假设在访问首页时列出所有上传的图片。

由于我们将客户端上传的图片全部保存在工程的 ./uploads 目录下，所以程序中应该有个名叫 listHandler() 的方法，用于在网页上列出该目录下存放的所有文件。暂时我们不考虑以缩略图的形式列出所有已上传图片，只需列出可供访问的文件名称即可。下面我们就来实现这个 listHandler() 方法：

```go
// 图片列表
func listHandler(writer http.ResponseWriter, request *http.Request) {
	// 读取文件夹
	infos, err := ioutil.ReadDir(UploadDir)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var html string
	for _, info := range infos {
		name := info.Name()
		html += `<li><a href="/view?id=` + name + `">` + name + `</a></li>`
	}

	// 输出图片列表
	io.WriteString(writer, `<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <title>Title</title>
		</head>
		<body>
			<ol>`+html+`</ol>
		</body>
		</html>
	`)

	return
}
```

从上面的 listHandler() 方法中可以看到，程序先从 ./uploads 目录中遍历得到所有文件并赋值到 infos 变量里。infos 是一个数组，其中的每一个元素都是一个文件对象。

然后，程序遍历 infos 数组并从中得到图片的名称，用于在后续的 HTML 片段中显示文件名和传入的参数内容。html 变量用于在 for 循序中将图片名称一一串联起来生成一段 HTML，最后调用 io.WriteString() 方法将这段 HTML 输出返回给客户端。

然后在 photos.go 程序的 main() 方法中，我们将对首页的访问映射到 listHandler() 方法。main() 方法的代码如下：


```go
func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/", listHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

这样在访问网站首页的时候，即可看到已上传的所有图片列表了。

![图片列表](https://lucklit.oss-cn-beijing.aliyuncs.com/written/Snip20191121_2.png)

不过，你是否注意到一个事实，我们在 photos.go 程序的 uploadHandler() 和 listHandler() 方法中都使用 io.WriteString() 方法输出 HTML。

正如你想到的那样，在业务逻辑处理程序中混杂 HTML 可不是什么好事情，代码多起来后会导致程序不够清晰，而且改动程序里边的 HTML 文本时，每次都要重新编译整个工程的源代码才能看到修改后的效果。

正确的做法是，应该将业务逻辑程序和表现层分离开来，各自单独处理。这时候，就需要使用网页模板技术了。

Go 标准库中的 html/template 包对网页模板有着良好的支持。接下来，让我们来了解如何在 photoweb.go 程序中用上 Go 的模板功能。

#### 渲染网页模板

