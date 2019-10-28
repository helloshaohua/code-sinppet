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
