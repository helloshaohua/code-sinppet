package main

import (
	"bufio"
	"bytes"
	"fmt"
)

func main() {
	wr := bytes.NewBuffer(nil)
	w := bufio.NewWriter(wr)
	data := []byte("Hello World")
	w.Write(data)
	fmt.Printf("未执行 Flush 缓冲区输出 %q\n")
}
