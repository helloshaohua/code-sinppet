### Go语言词频统计


从数据挖掘到语言学习本身，文本分析功能的应用非常广泛，这一节我们来分析一个例子，它是文本分析最基本的一种形式：统计出一个文件里单词出现的频率。

示例中频率统计后的结果可以以两种不同的方式显示，一种是将单词按照字母顺序把单词和频率排列出来，另一种是按照有序列表的方式把频率和对应的单词显示出来，完整的示例代码如下所示：

```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "-help" {
		fmt.Printf("usage: %s <file1> [<file2> [... <fileN>]]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	// 使用 make 初始化内置数据结构，用来保存从文件读到的每一个单词和对应的频率
	frequencyForWord := make(map[string]int)
	// 遍历从命令行得到的每一个文件，分析每一个文件后更新 frequencyForWord 的数据
	for _, filename := range commandLineFiles(os.Args[1:]) {
		updateFrequencies(filename, frequencyForWord)
	}
	reportByWords(frequencyForWord)
	wordsForFrequency := invertStringIntMap(frequencyForWord)
	reportByFrequency(wordsForFrequency)
}

// 实现跨平台获取命令行文件的处理
func commandLineFiles(files []string) []string {
	if runtime.GOOS == "windows" {
		args := make([]string, 0, len(files))
		for _, filename := range files {
			if matches, err := filepath.Glob(filename); err != nil {
				args = append(args, filename)
			} else if matches != nil {
				args = append(args, matches...)
			}
		}
		return args
	}
	return files
}

// 分析每一个文件后更新 frequencyForWord 的数据
func updateFrequencies(filename string, frequencyForWord map[string]int) {
	var file *os.File
	var err error
	if file, err = os.Open(filename); err != nil {
		log.Printf("os.Open error: %s\n", err)
		return
	}
	defer file.Close()
	readAndUpdateFrequencies(bufio.NewReader(file), frequencyForWord)
}

// 读取文件并更新 frequencyForWord 的数据
func readAndUpdateFrequencies(reader *bufio.Reader, frequencyForWord map[string]int) {
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("reader.ReadString error: %s\n", err)
			}
			break
		}

		for _, word := range SplitOnNonLetters(strings.TrimSpace(s)) {
			if len(word) > utf8.UTFMax || utf8.RuneCountInString(word) > 1 {
				frequencyForWord[strings.ToLower(word)] += 1
			}
		}
	}
}

// 用来在非单词字符上对一个字符串进行切分
func SplitOnNonLetters(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}

// 打印frequencyForWord数据出来
func reportByWords(frequencyForWord map[string]int) {
	words := make([]string, 0, len(frequencyForWord))
	wordWidth, frequencyWidth := 0, 0
	for word, frequency := range frequencyForWord {
		words = append(words, word)
		if width := utf8.RuneCountInString(word); width > wordWidth {
			wordWidth = width
		}

		if width := len(fmt.Sprint(frequency)); width > frequencyWidth {
			frequencyWidth = width
		}
	}

	// 对切片排序
	sort.Strings(words)

	gap := wordWidth + frequencyWidth - len("Word") - len("Frequency")

	fmt.Printf("Word %*s%s\n", gap, " ", "Frequency")

	for _, word := range words {
		fmt.Printf("%-*s %*d\n", wordWidth, word, frequencyWidth, frequencyForWord[word])
	}
}

// 反转统计结果
func invertStringIntMap(intForString map[string]int) map[int][]string {
	stringForInt := make(map[int][]string, len(intForString))
	for key, value := range intForString {
		stringForInt[value] = append(stringForInt[value], key)
	}
	return stringForInt
}

// 输出反转统计结果
func reportByFrequency(wordsForFrequency map[int][]string) {
	frequencies := make([]int, 0, len(wordsForFrequency))
	for frequency := range wordsForFrequency {
		frequencies = append(frequencies, frequency)
	}
	// 整型切片排序
	sort.Ints(frequencies)
	width := len(fmt.Sprint(frequencies[len(frequencies)-1]))
	fmt.Println("Frequency -> words")
	for _, frequency := range frequencies {
		words := wordsForFrequency[frequency]
		sort.Strings(words)
		fmt.Printf("%*d %s\n", width, frequency, strings.Join(words, ", "))
	}
}
```

程序的两种生成方式，如下所示。

```shell
go run test.go words.txt
Word   Frequency
among          1
and            4
api            1
automatically  1
based          1
be             2
became         2
can            2
client         1
clients        1
community      1
descriptions   1
developer      1
ever           1
focus          1
for            2
from           1
generated      1
get            1
go             1
great          1
grpc           4
http           1
implementation 1
in             1
introduced     1
is             1
it             1
its            2
languages      1
line           1
message        1
of             1
on             2
performance    1
polyglot       1
popular        2
protobuf       1
really         1
reason         1
separate       1
server         1
service        1
set            1
since          1
single         1
so             1
support        1
the            2
tool           1
use            1
was            2
without        1
writing        1
written        1
Frequency -> words
1 among, api, automatically, based, client, clients, community, descriptions, developer, ever, focus, from, generated, get, go, great, http, implementation, in, introduced, is, it, languages, line, message, of, performance, polyglot, protobuf, really, reason, separate, server, service, set, since, single, so, support, tool, use, without, writing, written
2 be, became, can, for, its, on, popular, the, was
4 and, grpc
```

即使是很小的文件，不同单词的数量也会非常大，所以这里只截取了部分结果。

第一种输出是比较直接的，我们可以使用一个 map [string] int 类型的结构来保存每一个单词的频率，但是要得到第二种输出结果我们需要将整个映射反转成多值类型的映射，如 map[int][] string，也就是说，键是频率而值则是所有具有这个频率的单词。

接下来我们将从程序的 main() 函数开始，从上到下分析。

```go
func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "-help" {
		fmt.Printf("usage: %s <file1> [<file2> [... <fileN>]]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	// 使用 make 初始化内置数据结构，用来保存从文件读到的每一个单词和对应的频率
	frequencyForWord := make(map[string]int)
	// 遍历从命令行得到的每一个文件，分析每一个文件后更新 frequencyForWord 的数据
	for _, filename := range commandLineFiles(os.Args[1:]) {
		updateFrequencies(filename, frequencyForWord)
	}
	reportByWords(frequencyForWord)
	wordsForFrequency := invertStringIntMap(frequencyForWord)
	reportByFrequency(wordsForFrequency)
}
```

main() 函数首先分析命令行参数，之后再进行相应处理。

我们使用复合语法创建一个空的映射，用来保存从文件读到的每一个单词和对应的频率，接着我们遍历从命令行得到的每一个文件，分析每一个文件后更新 frequencyForWord 的数据。

得到第一个映射之后，我们就可以输出第一个报告了，一个按照字母表顺序排序的单词列表和对应的出现频率，然后我们创建一个反转的映射，输出第二个报告，一个按出现频率排序的列表。

```go
// 实现跨平台获取命令行文件的处理
func commandLineFiles(files []string) []string {
	if runtime.GOOS == "windows" {
		args := make([]string, 0, len(files))
		for _, filename := range files {
			if matches, err := filepath.Glob(filename); err != nil {
				args = append(args, filename)
			} else if matches != nil {
				args = append(args, matches...)
			}
		}
		return args
	}
	return files
}
```

因为 Unix 类系统（如 Linux 或 Mac OS X 等）的 shell 默认会自动处理通配符（也就是说，*.txt 能匹配任意后缀为 .txt 的文件，如 README.txt 和 INSTALL.txt 等），而 Windows 平台的 shell 程序（cmd.exe）不支持通配符，所以如果用户在命令行输入 *.txt，那么程序只能接收到 *.txt。

为了保持平台之间的一致性，这里使用 commandLineFiles() 函数来实现跨平台的处理，当程序运行在 Windows 平台时，把文件名通配功能给实现了。

```go
// 分析每一个文件后更新 frequencyForWord 的数据
func updateFrequencies(filename string, frequencyForWord map[string]int) {
	var file *os.File
	var err error
	if file, err = os.Open(filename); err != nil {
		log.Printf("os.Open error: %s\n", err)
		return
	}
	defer file.Close()
	readAndUpdateFrequencies(bufio.NewReader(file), frequencyForWord)
}
```

updateFrequencies() 函数纯粹就是用来处理文件的，它打开给定的文件，并使用 defer 让函数返回时关闭文件句柄，这里我们将文件作为一个 *bufio.Reader（使用 bufio.NewReader() 函数创建）传给 readAndUpdateFrequencies() 函数，因为这个函数是以字符串的形式一行一行地读取数据的，所以实际的工作都是在 readAndUpdateFrequencies() 函数里完成的，代码如下。

```go
// 读取文件并更新 frequencyForWord 的数据
func readAndUpdateFrequencies(reader *bufio.Reader, frequencyForWord map[string]int) {
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("reader.ReadString error: %s\n", err)
			}
			break
		}

		for _, word := range SplitOnNonLetters(strings.TrimSpace(s)) {
			if len(word) > utf8.UTFMax || utf8.RuneCountInString(word) > 1 {
				frequencyForWord[strings.ToLower(word)] += 1
			}
		}
	}
}
```

第一部分的代码我们应该很熟悉了，用了一个无限循环来一行一行地读一个文件，当读到文件结尾或者出现错误（这种情况下需要将错误报告给用户）的时候就退出循环，但并不退出程序，因为还有很多其他的文件需要去处理。

任意一行都可能包括标点、数字、符号或者其他非单词字符，所以我们需要逐个单词地去读，将每一行分隔成若干个单词并使用 SplitOnNonLetters() 函数忽略掉非单词的字符，并且过滤掉字符串开头和结尾的空白。

只需要记录含有两个以上（包括两个）字母的单词，可以通过使用 if 语句，如 utf8.RuneCountlnString(word) > 1 来完成。

上面描述的 if 语句有一点性能损耗，因为它会分析整个单词，所以在这个程序里我们增加了一个判断条件，用来检査这个单词的字节数是否大于 utf8.UTFMax（utf8.UTFMax 是一个常量，值为 4，用来表示一个 UTF-8 字符最多需要几个字节）。

```go
// 用来在非单词字符上对一个字符串进行切分
func SplitOnNonLetters(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}
```

SplitOnNonLetters() 函数用来在非单词字符上对一个字符串进行切分，首先我们为 strings.FieldsFunc() 函数创建一个匿名函数 notALetter，如果传入的是字符那就返回 false，否则返回 true，然后返回调用函数 strings.FieldsFunc() 的结果，调用的时候将给定的字符串和 notALetter 作为它的参数。

```go
// 打印frequencyForWord数据出来
func reportByWords(frequencyForWord map[string]int) {
	words := make([]string, 0, len(frequencyForWord))
	wordWidth, frequencyWidth := 0, 0
	for word, frequency := range frequencyForWord {
		words = append(words, word)
		if width := utf8.RuneCountInString(word); width > wordWidth {
			wordWidth = width
		}

		if width := len(fmt.Sprint(frequency)); width > frequencyWidth {
			frequencyWidth = width
		}
	}

	// 对切片排序
	sort.Strings(words)

	gap := wordWidth + frequencyWidth - len("Word") - len("Frequency")

	fmt.Printf("Word %*s%s\n", gap, " ", "Frequency")

	for _, word := range words {
		fmt.Printf("%-*s %*d\n", wordWidth, word, frequencyWidth, frequencyForWord[word])
	}
}
```

计算出了 frequencyForWord 之后，调用 reportByWords() 将它的数据打印出来，因为我们希望输出结果是按照字母顺序排序好的，所以首先要创建一个空的容量足够大的 []string 切片来保存所有在 frequencyForWord 里的单词。

第一个循环遍历映射里的所有项，把每个单词追加到 words 字符串切片里去，因为 words 的容量足够大了，所以 append() 函数只需要把给定的单词追加到第 len(words) 个索引位置上即可，words 的长度会自动增加 1。

得到了 words 切片之后，对它进行排序，这个在 readAndUpdateFrequencies() 函数中已经处理好了。

经过排序之后我们打印两列标题，第一个是 "Word"，为了能让 Frequency 最后一个字符 y 右对齐，需要在 "Word" 后打印一些空格，通过%*s可以实现的打印固定长度的空白，另一种办法是可以使用 %s来打印 strings.Repeat(" ", gap) 返回的字符串。

最后，我们将单词和它们的频率用两列方式按照字母顺序打印出来。

```go
func invertStringIntMap(intForString map[string]int) map[int][]string {
    stringsForInt := make(map[int][]string, len(intForString))
    for key, value := range intForString {
        stringsForInt[value] = append(stringsForInt[value], key)
    }
    return stringsForInt
}
```

上面的函数首先创建一个空的映射，用来保存反转的结果，但是我们并不知道它到底要保存多少个项，因此我们假设它和原来的映射容量一样大，然后简单地遍历原来的映射，将它的值作为键保存到反转的映射里，并将键增加到对应的值里去，新的映射的值就是一个字符串切片，即使原来的映射有多个键对应同一个值，也不会丢掉任何数据。


```go
// 输出反转统计结果
func reportByFrequency(wordsForFrequency map[int][]string) {
	frequencies := make([]int, 0, len(wordsForFrequency))
	for frequency := range wordsForFrequency {
		frequencies = append(frequencies, frequency)
	}
	// 整型切片排序
	sort.Ints(frequencies)
	width := len(fmt.Sprint(frequencies[len(frequencies)-1]))
	fmt.Println("Frequency -> words")
	for _, frequency := range frequencies {
		words := wordsForFrequency[frequency]
		sort.Strings(words)
		fmt.Printf("%*d %s\n", width, frequency, strings.Join(words, ", "))
	}
}
```

这个函数的结构和 reportByWords() 函数很相似，它首先创建一个切片用来保存频率，并按照频率升序排列，然后再计算需要容纳的最大长度并以此作为第一列的宽度，之后输出报告的标题，最后，遍历输出所有的频率并按照字母升序输出对应的单词，如果一个频率有超过两个对应的单词则单词之间使用逗号分隔开。

