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
