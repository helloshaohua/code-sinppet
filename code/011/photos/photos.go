package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/", listHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

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
