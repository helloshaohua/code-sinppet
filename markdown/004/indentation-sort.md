### Go语言缩进排序

本节将通过实例为大家演示如何将字符串按照等级（缩进级别）进行排序，完整代码如下所示。

```go
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
```

注意 SortedIndentedStrings() 函数有一个很重要的前提就是，字符串的缩进是通过读到的空格和缩进的个数来决定的，我们来看一下 main 函数和输出结果，为了方便对比，我们将排序前的结果放在左边，排序后的结果放在右边。

```go
func main() {
	fmt.Println("|     Original      |       Sorted      |")
	fmt.Println("|-------------------|-------------------|")
	sorted := sortedIndentedStrings(original)
	for i := range original {
		fmt.Printf("|%-19s|%-19s|\n", original[i], sorted[i])
	}
}
```

```shell
$ go run test.go
|     Original      |       Sorted      |
|-------------------|-------------------|
|Nonmetals          |Alkali Metals      |
|    Hydrogen       |    Lithium        |
|    Carbon         |    Potassium      |
|    Nitrogen       |    Sodium         |
|    Oxygen         |Inner Transitionals|
|Inner Transitionals|    Actinides      |
|    Lanthanides    |        Curium     |
|        Europium   |        Plutonium  |
|        Cerium     |        Uranium    |
|    Actinides      |    Lanthanides    |
|        Uranium    |        Cerium     |
|        Plutonium  |        Europium   |
|        Curium     |Nonmetals          |
|Alkali Metals      |    Carbon         |
|    Lithium        |    Hydrogen       |
|    Sodium         |    Nitrogen       |
|    Potassium      |    Oxygen         |
```

其中，SortedIndentedStrings() 函数和它的辅助函数使用到了递归、函数引用以及指向切片的指针等。

```go
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
```

sort.Interface 接口定义了 3 个方法 Len()、Less() 和 Swap()，它们的函数签名和 Entries 中的同名方法是一样的，这就意味着我们可以使用标准库里的 sort.Sort() 函数来对一个 Entries 进行排序。

```go
// 排序缩进字符串
func sortedIndentedStrings(slice []string) []string {
	entries := populateEntries(slice)
	return sortedEntries(entries)
}
```

导出的函数 SortedIndentedStrings() 就做了这个工作，虽然我们已经对它进行了重构，让它把所有东西都传递给辅助函数，函数 populateEntries() 传入一个 []string 并返回一个对应的 Entries（[]Entry 类型）。

而函数 sortedEntries() 需要传入一个 Entries，然后返回一个排过序的 []string（根据缩进的级别进行排序）。

```go
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
```

populateEntries() 函数首先以字符串的形式得到给定切片里的一级缩进（如有 4 个空格的字符串）和它占用的字节数，然后创建一个空的 Entries，并遍历切片里的每一个字符串，判断该字符串的缩进级别，再创建一个用于排序的键。

下一步，调用自定义函数 addEntry()，将当前字符串的级别、键、字符串本身，以及指向 entries 的地址作为参数，这样 addEntry() 就能创建一个新的 Entry 并能够正确地将它追加到 entries 里去，最后返回 entries。

```go
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
```

computeIndent() 函数主要是用来判断缩进使用的是什么字符，例如空格或者缩进符等，以及一个缩进级别占用多少个这样的字符。

因为第一级的字符串可能没有缩进，所以函数必须迭代所有的字符串，一旦它发现某个字符串的行首是空格或者缩进，函数马上返回表示缩进的字符以及一个缩进所占用的字符数。

```go
// 条目分组细化
func addEntry(level int, key string, value string, entries *Entries) {
	if level == 0 {
		*entries = append(*entries, Entry{key: key, value: value, children: make(Entries, 0)})
	} else {
		addEntry(level-1, key, value, &((*entries)[entries.Len()-1].children))
	}
}
```

addEntry() 是一个递归函数，它创建一个新的 Entry，如果这个 Entry 的 level 是 0，那就直接增加到 entries 里去，否则，就将它作为另一个 Entry 的子集。

我们必须确定这个函数传入的是一个 *Entries 而不是传递一个 entries 引用（切片的默认行为），因为我们是要将数据追加到 entries 里，追加到一个引用会导致无用的本地副本且原来的数据实际上并没有被修改。

如果 level 是 0，表明这个字符串是顶级项，因此必须将它直接追加到 *entries，实际上情况要更复杂一些，因为 level 是相对传入的 *entries 而言的，第一次调用 addEntry() 时，*entries 是一个第一级的 Entries，但函数进入递归后，*entries 就可能是某个 Entry 的子集。

我们使用内置的 append() 函数来追加新的 Entry，并使用 * 操作符获得 entries 指针指向的值，这就保证了任何改变对调用者来说都是可见的，新增的 Entry 包含给定的 key 和 value，以及一个空的子 Entries，这是递归的结束条件。

如果 level 大于 0，则我们必须将它追加到上一级 Entry 的 children 字段里去，这里我们只是简单地递归调用 addEntry() 函数，最后一个参数可能是我们目前为止见到的最复杂的表达式了。

子表达式 entries.Len() - 1 产生一个 int 型整数，表示 *entries 指向的 Entries 值的最后一个条目的索引位置（注意 Entries.Len() 传入的是一个 Entries 值而不是 *Entries 指针，不过Go语言也可以自动对 entries 指针进行解引用并调用相应的方法）。

完整的表达式（&(...) 除外）访问了 Entries 最后一个 Entry 的 children 字段（这也是一个 Entries 类型），所以如果把这个表达式作为一个整体，实际上我们是将 Entries 里最后一个 Entry 的 children 字段的内存地址作为递归调用的参数，因为 addEntry() 最后一个参数是 *Entries 类型的。

为了帮助大家弄清楚到底发生了什么，下面的代码和上述代码中 else 代码块中的那个调用是一样的。

```go
theEntries := *entries
lastEntry := &theEntries[theEntries.Len()-1]
addEntry(level-1, key, value, &lastEntry.children)
```

首先，我们创建 theEntries 变量用来保存 *entries 指针指向的值，这里没有什么开销因为不会产生复制，实际上 theEntries 相当于一个指向 Entries 值的别名。

然后我们取得最后一项的内存地址（即一个指针），如果不取地址的话就会取到最后一项的副本，最后递归调用 addEntry() 函数，并将最后一项的 children 字段的地址作为参数传递给它。

```go
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
```

当调用 sortedEntries() 函数的时候，Entries 显示的结构和原先程序输出的字符串是一样的，每一个缩进的字符串都是上一级缩进的子级，而且还可能有下一级的缩进，依次类推。

创建了 Entries 之后，SortedIndentedStrings() 函数调用上面这个函数去生成一个排好序的字符串切片 []string，这个函数首先创建一个空的 []string 用来保存最后的结果，然后对 entries 进行排序。

Entries 实现了 sort.Interface 接口，因此我们可以直接使用 sort.Sort() 函数根据 Entry 的 key 字段来对 Entries 进行排序（这是 Entries.Less() 的实现方式），这个排序只是作用于第一级的 Entry，对其他未排序的子集是没有任何影响的。

为了能够对 children 字段以及 children 的 children 等进行递归排序，函数遍历第一级的每一个项并调用 populateIndentedStrings() 函数，传入这个 Entry 类型的项和一个指向 []string 切片的指针。

切片可以传递给函数并由函数更新内容（如替换切片里的某些项），但是这里需要往切片里新增一些数据，所以这里将一个指向切片的指针（也就是指针的指针）作为参数传进去，并将指针指向的内容设置为 append() 函数的返回结果，可能是一个新的切片，也可能是原先的切片。

另一种办法就是传入切片的值，然后返回 append() 之后的切片，但是必须将返回的结果赋值给原来的切片变量（例如 slice = function(slice)），不过这么做的话，很难正确地使用递归函数。

```go
// 对子条目集合排序
func populateIndentedStrings(entry Entry, indentedSlice *[]string) {
	// 将排序后的条目放到目标结果集中
	*indentedSlice = append(*indentedSlice, entry.value)
	sort.Sort(entry.children)
	for _, child := range entry.children {
		populateIndentedStrings(child, indentedSlice)
	}
}
```

populateIndentedStrings() 函数将顶级项追加到创建的切片，然后对顶级项的子项进行排序，并递归调用自身对每一个子项做同样的处理，这就相当于对每一项的子项以及子项的子项等都做了排序，所以整个字符串切片就是已经排好序的了。