package main

import (
	"fmt"
	"sort"
	"strings"
)

// 数据源
var original = []string{
	"Nonmetals",
	"    Hydrogen",
	"    Carbon",
	"    Nitrogen",
	"    Oxygen",
	"Inner Transitionals",
	"    Lanthanides",
	"        Europium",
	"        Cerium",
	"    Actinides",
	"        Uranium",
	"        Plutonium",
	"        Curium",
	"Alkali Metals",
	"    Lithium",
	"    Sodium",
	"    Potassium",
}

// 条目结构体定义
type Entry struct {
	key      string
	value    string
	children Entries
}

// 条目结构体集合定义
type Entries []Entry

// 获取元素个数
func (entries Entries) Len() int {
	return len(entries)
}

// 元素与元素之间的键比较大小
func (entries Entries) Less(i, j int) bool {
	return entries[i].key < entries[j].key
}

// 元素位置交换
func (entries Entries) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

func main() {
	fmt.Println("|     Original      |       Sorted      |")
	fmt.Println("|-------------------|-------------------|")
	sorted := sortedIndentedStrings(original)
	for i := range original {
		fmt.Printf("|%-19s|%-19s|\n", original[i], sorted[i])
	}
}

// 排序缩进字符串
func sortedIndentedStrings(slice []string) []string {
	entries := populateEntries(slice)
	return sortedEntries(entries)
}

// 填充排序条目
func populateEntries(slice []string) Entries {
	indent, indentSize := computeIndent(slice)
	entries := make(Entries, 0)
	for _, item := range slice {
		i, level := 0, 0
		for strings.HasPrefix(item[i:], indent) {
			i += indentSize
			level++
		}
		key := strings.ToLower(strings.TrimSpace(item))
		addEntry(level, key, item, &entries)
	}
	return entries
}

// 计算缩进字符及缩进位置
func computeIndent(slice []string) (string, int) {
	for _, item := range slice {
		if len(item) > 0 && (item[0] == ' ' || item[0] == '\t') {
			whitespace := rune(item[0])
			for i, char := range item[1:] {
				if char != whitespace {
					i++
					return strings.Repeat(string(whitespace), i), i
				}
			}
		}
	}
	return "", 0
}

// 条目分组细化
func addEntry(level int, key string, value string, entries *Entries) {
	if level == 0 {
		*entries = append(*entries, Entry{key: key, value: value, children: make(Entries, 0)})
	} else {
		addEntry(level-1, key, value, &((*entries)[entries.Len()-1].children))
	}
}

// 条目集合排序
func sortedEntries(entries Entries) []string {
	indentedSlice := make([]string, 0)
	// 第一级的条目排序
	sort.Sort(entries)
	for _, entry := range entries {
		// 对子条目集合排序
		populateIndentedStrings(entry, &indentedSlice)
	}
	return indentedSlice
}

// 对子条目集合排序
func populateIndentedStrings(entry Entry, indentedSlice *[]string) {
	// 将排序后的条目放到目标结果集中
	*indentedSlice = append(*indentedSlice, entry.value)
	sort.Sort(entry.children)
	for _, child := range entry.children {
		populateIndentedStrings(child, indentedSlice)
	}
}
