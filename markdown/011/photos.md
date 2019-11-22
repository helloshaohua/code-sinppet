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

使用 Go 标准库提供的 html/template 包，可以让我们将 HTML 从业务逻辑程序中抽离出来形成独立的模板文件，这样业务逻辑程序只负责处理业务逻辑部分和提供模板需要的数据，模板文件负责数据要表现的具体形式。

然后模板解析器将这些数据以定义好的模板规则结合模板文件进行渲染，最终将渲染后的结果一并输出，构成一个完整的网页。

下面我们把 photos.go 程序的 uploadHandler() 和 listHandler() 方法中的 HTML 文本 抽出，生成模板文件。

新建一个 views 目录用于存放 html 模板文件。



在 views 目录下新建一个名为 upload.html 的文件，内容如下：

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Upload</title>
</head>
<body>
<form method="post" action="/upload" enctype="multipart/form-data">
    Choose an image to upload:
    <input name="image" type="file">
    <input value="Upload" type="submit">
</form>
</body>
</html>
```

然后在 views 目录下 新建一个名为 list.html 的文件，内容如下：

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
    <ol>
        {{range $.images}}
            <li><a href="/view?id{{.|urlquery}}">{{.|html}}</a></li>
        {{end}}
    </ol>
</body>
</html>
```

在上述模板中，双大括号 {{}} 是区分模板代码和 HTML 的分隔符，括号里边可以是要显示输出的数据，或者是控制语句，比如 if 判断式或者 range 循环体等。

range 语句在模板中是一个循环过程体，紧跟在 range 后面的必须是一个 array、slice 或 map 类型的变量。在 list.html 模板中，images 是一组 string 类型的切片。

在使用 range 语句遍历的过程中，. 即表示该循环体中的当前元素，.|formatter 表示对当前这个元素的值以 formatter 方式进行格式化输出，比如 .|urlquery} 即表示对当前元素的值进行转换以适合作为 URL 一部分，而 {{.|html 表示对当前元素的值进行适合用于 HTML 显示的字符转化，比如">"会被转义成"&gt;"。

如果 range 关键字后面紧跟的是 map 这样的多维复合结构，循环体中的当前元素可以用 .key1.key2.keyN 这样的形式表示。

如果要更改模板中默认的分隔符，可以使用 template 包提供的 Delims() 方法。

在了解模板语法后，接着我们修改 photos.go 源文件，引入 html/template 包，并修改 uploadHandler() 和 listHandler() 方法，具体如下所示。

> 修改uploadHandler()

```go
// 图片上传
func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	// 显示表单
	if request.Method == "GET" {
		// 解析模板文件
		tpl, err := template.ParseFiles(getFilePath(TemplateDir, "upload.html"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// 将模板写入响应对象
		err = tpl.Execute(writer, nil)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回停止处理
		return
	}

	// 上传图片
	if request.Method == "POST" {
		// ......
	}
}
```

> 修改listHandler()

```go
// 图片列表
func listHandler(writer http.ResponseWriter, request *http.Request) {
	// 读取文件夹
	infos, err := ioutil.ReadDir(UploadDir)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取图片文件名将其放入到一个切片
	images := make([]string, 0)
	for _, info := range infos {
		images = append(images, info.Name())
	}

	// 解析模板文件
	tpl, err := template.ParseFiles(getFilePath(TemplateDir, "list.html"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将模板文件写入响应对应
	err = tpl.Execute(writer, map[string]interface{}{"images": images})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
```

在这个过程中我们还定义了一个表示 views 视图文件目录的常量：

```go
const (
	TemplateDir = "./code/011/photos/views"
)
```

当我们引入了 Go 标准库中的 html/template 包，实现了业务逻辑层与表现层分离后(你可以通过浏览器验证)，对模板渲染逻辑去重，编写并使用通用模板渲染方法 renderHtml()，这让业务逻辑处理层的代码看起来确实要清晰简洁许多。

不过，直觉敏锐的你可能已经发现，无论是重构后的 uploadHandler() 还是 listHandler() 方法，每次调用这两个方法时都会重新读取并渲染模板。很明显，这很低效，也比较浪费资源，有没有一种办法可以让模板只加载一次呢？

答案是肯定的，聪明的你可能已经想到怎么对模板进行缓存了。

#### 模板缓存

对模板进行缓存，即指一次性预加载模板。我们可以在 photos.go 程序初始化运行的时候，将所有模板一次性加载到程序中。正好 Go 的包加载机制允许我们在 init() 函数中做这样的事情，init() 会在 main() 函数之前执行。

首先，我们在 photos.go 程序中声明并初始化一个全局变量 templates，用于存放所有模板内容：

```go
// 全局变量 templates 用于存放所有模板内容
var templates = make(map[string]*template.Template)
```

templates 是一个 map 类型的复合结构，map 的键（key）是字符串类型，即模板的名字，值（value）是 *template.Template 类型。

接着，我们在 photos.go 程序的 init() 函数中一次性加载所有模板：

```go
// 全局变量 templates 用于存放所有模板内容
var templates = make(map[string]*template.Template)

// 初始化函数完成初始化工作
func init() {
	for _, filename := range []string{"list", "upload"} {
		tpl := template.Must(template.ParseFiles(getFilePath(TemplateDir, filename+".html")))
		templates[filename] = tpl
	}
}
```

在上面的代码中，我们在 template.ParseFiles() 方法的外层强制使用 template.Must() 进行封装，template.Must() 确保了模板不能解析成功时，一定会触发错误处理流程。之所以这么做，是因为倘若模板不能成功加载，程序能做的唯一有意义的事情就是退出。

在 range 语句中，包含了我们希望加载的 upload.html 和 list.html 两个模板，如果我们想加载更多模板，只需往这个数组中添加更多元素即可。当然，最好的办法应该是将 /views 模板文件夹进行遍历和预加载。

如果需要加载新的模板，只需在这个 /views 文件夹中新建模板即可。这样做的好处是不用反复修改代码即可重新编译程序，而且实现了业务层和表现层真正意义上的分离。

接下来适当地对 init() 方法中的代码进行改写，好让程序初始化时即可预加载该目录下的所有模板文件，如下列代码所示：

```go
// 全局变量 templates 用于存放所有模板内容
var templates = make(map[string]*template.Template)

// 初始化函数完成初始化工作
func init() {
	// 读取文件夹
	infos, err := ioutil.ReadDir(TemplateDir)
	if err != nil {
		panic(err)
		return
	}

	// 遍历文件数组
	for _, info := range infos {
		filename := info.Name()
		if ext := path.Ext(filename); ext != ".html" {
			continue
		}

		// 解析HTML模板文件
		tpl := template.Must(template.ParseFiles(getFilePath(TemplateDir, filename)))

		// 添加到全部变量
		templates[filename] = tpl
	}
}
```

添加 renderHTML 辅助函数，用于渲染预加载的 HTML 模板文件，定义如下：

```go
// 渲染HTML模板
func renderHTML(writer http.ResponseWriter, filename string, data map[string]interface{}) error {
	if tpl, ok := templates[filename]; ok {
		return tpl.Execute(writer, data)
	}
	return errors.New(fmt.Sprintf("%s file is not exists!", filename))
}
```

此时， 我们可以将 uploadHandler 和 listHandler 处理函数进行调整，将渲染 HTML 的方式进行调整，避免每次调整处理函数时都进行 HTML 的加载过程，调整代码如下：

> uploadHandler处理函数

```go
// 图片上传
func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	// 显示表单
	if request.Method == "GET" {
		err := renderHTML(writer, "upload.html", nil)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回停止处理
		return
	}

	// 上传图片
	if request.Method == "POST" {
		// ......
	}
}
```

> listHandler处理函数

```go
// 图片列表
func listHandler(writer http.ResponseWriter, request *http.Request) {
	// 读取文件夹
	infos, err := ioutil.ReadDir(UploadDir)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取图片文件名将其放入到一个切片
	images := make([]string, 0)
	for _, info := range infos {
		images = append(images, info.Name())
	}

	// 渲染 HTML
	err = renderHTML(writer, "list.html", map[string]interface{}{"images": images})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
```

做完这些工作可以进行浏览器进行验证是否达到预览或是否有错误产生！

#### 错误处理

在前面的代码中，有不少地方对于出错处理都是直接返回 http.Error() 50x 系列的服务端内部错误。从 DRY 的原则来看，不应该在程序中到处使用一样的代码。我们可以定义一个名为 check() 的方法，用于统一捕获 50x 系列的服务端内部错误：

```go
func check(err error) {
	if err != nil {
		panic(err)
	}
}
```

此时，我们可以将 photos.go 程序中出现的以下代码：

```go
if err != nil {
    http.Error(w, err.Error(),http.StatusInternalServerError)
    return
}
```

统一替换为 check() 处理：

```go
check(err)
```

错误处理虽然简单很多，但是也带来一个问题。由于发生错误触发错误处理流程必然会引发程序停止运行，这种改法有点像搬起石头砸自己的脚。

其实我们可以换一种思维方式。尽管我们从书写上能保证大多数错误都能得到相应的处理，但根据墨菲定律，有可能出问题的地方就一定会出问题，在计算机程序里尤其如此。如果程序中我们正确地处理了 99 个错误，但若有一个系统错误意外导致程序出现异常，那么程序同样还是会终止运行。

我们不能预计一个工程里边会出现多少意外的情况，但是不管什么意外，只要会触发错误处理流程，我们就有办法对其进行处理。如果这样思考，那么前面这种改法又何尝不是置死地而后生呢？

接下来，让我们了解如何处理 panic 导致程序崩溃的情况。

#### 巧用闭包避免程序运行时出错崩溃

Go 支持闭包。闭包可以是一个函数里边返回的另一个匿名函数，该匿名函数包含了定义在它外面的值。使用闭包，可以让我们网站的业务逻辑处理程序更安全地运行。

我们可以在 photoweb 程序中针对所有的业务逻辑处理函数（listHandler()、viewHandler() 和 uploadHandler()）再进行一次包装。

在如下的代码中，我们定义了一个名为 safeHandler() 的函数，该函数有一个参数并且返回一个值，传入的参数和返回值都是一个函数，且都是http.HandlerFunc类型，这种类型的函数有两个参数：http.ResponseWriter 和 *http.Request。

函数规格同 photoweb 的业务逻辑处理函数完全一致。事实上，我们正是要把业务逻辑处理函数作为参数传入到 safeHandler() 方法中，这样任何一个错误处理流程向上回溯的时候，我们都能对其进行拦截处理，从而也能避免程序停止运行：

```go
// 统一处理错误
func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				// 或者输出自定义的 50x 错误页面
				// writer.WriteHeader(http.StatusInternalServerError)
				// renderHTML(writer, "error", map[string]interface{}{"error": err.Error()})

				// logging
				fcn := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
				log.Printf("WARN: panic in %v - %v\n", fcn, err)
				log.Println(string(debug.Stack()))
			}
		}()

		// 执行处理函数
		fn(writer, request)
	}
}
```

在上述这段代码中，我们巧妙地使用了 defer 关键字搭配 recover() 方法终结 panic 的肆行。safeHandler() 接收一个业务逻辑处理函数作为参数，同时调用这个业务逻辑处理函数。该业务逻辑函数执行完毕后，safeHandler() 中 defer 指定的匿名函数会执行。

倘若业务逻辑处理函数里边引发了 panic，则调用 recover() 对其进行检测，若为一般性的错误，则输出 HTTP 50x 出错信息并记录日志，而程序将继续良好运行。

要应用 safeHandler() 函数，只需在 main() 中对各个业务逻辑处理函数做一次包装，如下面的代码所示：

```go
func main() {
	http.HandleFunc("/upload", safeHandler(uploadHandler))
	http.HandleFunc("/view", safeHandler(viewHandler))
	http.HandleFunc("/", safeHandler(listHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

#### 动态请求和静态资源分离

你一定还有一个疑问，那就是前面的业务逻辑层都是动态请求，但若是针对静态资源（比如 CSS 和 JavaScript 等），是没有业务逻辑处理的，只需提供静态输出。在 Go 里边，这当然是可行的。

还记得前面我们在 viewHandler() 函数里边有用到 http.ServeFile() 这个方法吗？net/http 包提供的这个 ServeFile() 函数可以将服务端的一个文件内容读写到 http.Response-Writer 并返回给请求来源的 *http.Request 客户端。

用前面介绍的闭包技巧结合这个 http.ServeFile() 方法，我们就能轻而易举地实现业务逻辑的动态请求和静态资源的完全分离。

假设我们有 ./public 这样一个存放 css/、js/、images/ 等静态资源的目录，原则上所有如下的请求规则都指向该 ./public 目录下相对应的文件：

```text
[GET] /assets/css/*.css
[GET] /assets/js/*.js
[GET] /assets/images/*.js
```

然后，我们定义一个名为 staticDirHandler() 的方法，用于实现上述需求：

```go
// 静态资源处理
func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(writer http.ResponseWriter, request *http.Request) {
		f := staticDir + request.URL.Path[len(prefix)-1:]
		if flags&ListDir == 0 {
			if exists := isExists(f); !exists {
				http.NotFound(writer, request)
				return
			}
		}
		http.ServeFile(writer, request, f)
	})
}
```

如此即完美实现了静态资源和动态请求的分离。

当然，我们要思考是否确实需要用 Go 来提供静态资源的访问。如果使用外部 Web 服务器（比如 Nginx 等），就没必要使用 Go 编写的静态文件服务了。在本机做开发时有一个程序内置的静态文件服务器还是很实用的。

#### 重构

经过前面对 photoweb 程序一一重整之后，整个工程的目录结构如下：

```go
photos
├── photos.go
├── public
│   ├── css
│   ├── images
│   └── js
├── test
│   └── test.go
├── uploads
└── views
    ├── list.html
    └── upload.html
```

photos.go 程序的源码最终如下所示。

```go
package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	UploadDir   = "./code/011/photos/uploads"
	TemplateDir = "./code/011/photos/views"
	ListDir     = 0x0001
)

// 全局变量 templates 用于存放所有模板内容
var templates = make(map[string]*template.Template)

// 初始化函数完成初始化工作
func init() {
	// 读取文件夹
	infos, err := ioutil.ReadDir(TemplateDir)
	if err != nil {
		panic(err)
		return
	}

	// 遍历文件数组
	for _, info := range infos {
		filename := info.Name()
		if ext := path.Ext(filename); ext != ".html" {
			continue
		}

		// 解析HTML模板文件
		tpl := template.Must(template.ParseFiles(getFilePath(TemplateDir, filename)))

		// 添加到全部变量
		templates[filename] = tpl
	}
}

// 静态资源处理
func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(writer http.ResponseWriter, request *http.Request) {
		f := staticDir + request.URL.Path[len(prefix)-1:]
		if flags&ListDir == 0 {
			if exists := isExists(f); !exists {
				http.NotFound(writer, request)
				return
			}
		}
		http.ServeFile(writer, request, f)
	})
}

// 统一处理错误
func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				// 或者输出自定义的 50x 错误页面
				// writer.WriteHeader(http.StatusInternalServerError)
				// renderHTML(writer, "error", map[string]interface{}{"error": err.Error()})

				// logging
				fcn := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
				log.Printf("WARN: panic in %v - %v\n", fcn, err)
				log.Println(string(debug.Stack()))
			}
		}()

		// 执行处理函数
		fn(writer, request)
	}
}

// 图片列表
func listHandler(writer http.ResponseWriter, request *http.Request) {
	// 读取文件夹
	infos, err := ioutil.ReadDir(UploadDir)
	check(err)

	// 获取图片文件名将其放入到一个切片
	images := make([]string, 0)
	for _, info := range infos {
		images = append(images, info.Name())
	}

	// 渲染 HTML
	err = renderHTML(writer, "list.html", map[string]interface{}{"images": images})
	check(err)

	return
}

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
	check(err)

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

// 图片上传
func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	// 显示表单
	if request.Method == "GET" {
		// 渲染 HTML
		err := renderHTML(writer, "upload.html", nil)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回停止处理
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

// 渲染HTML模板
func renderHTML(writer http.ResponseWriter, filename string, data map[string]interface{}) error {
	if tpl, ok := templates[filename]; ok {
		return tpl.Execute(writer, data)
	}
	return errors.New(fmt.Sprintf("%s file is not exists!", filename))
}

// 错误检测
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// 获取文件路径
func getFilePath(dir, filename string) string {
	return dir + "/" + filename
}

// 获取新的上传文件名
func getNewFileNameForUpload(filename string) string {
	fs := strings.Split(filename, ".")
	ext := strings.ToLower(fs[len(fs)-1])
	t := time.Unix(time.Now().Unix(), 0).Format("2006-01-02-15-04-05")
	m := strconv.FormatInt(time.Now().UnixNano(), 10)
	return t + m + "." + ext
}

// 检测文件是否存在
func isExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func main() {
	mux := http.NewServeMux()
	staticDirHandler(mux, "/assets/", "./public", 0)
	mux.HandleFunc("/upload", safeHandler(uploadHandler))
	mux.HandleFunc("/view", safeHandler(viewHandler))
	mux.HandleFunc("/", safeHandler(listHandler))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
```

#### 更多资源

Go 的第三方库很丰富，无论是对于关系型数据库驱动还是非关系型的键值存储系统的接入，都有着良好的支持，而且还有丰富的 Go语言 Web 开发框架以及用于 Web 开发的相关工具包。可以访问 [http://godashboard.appspot.com/project](http://godashboard.appspot.com/project)，了解更多第三方库的详细信息。
