package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host\n", os.Args[0])
		os.Exit(1)
	}

	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	addr, _ := net.LookupHost("baidu.com")
	fmt.Printf("%+v\n", addr)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	ip := net.ParseIP(tcpAddr.IP.String())
	fmt.Println(ip.DefaultMask())

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)

	bytes, err := ioutil.ReadAll(conn)
	checkError(err)

	fmt.Println(string(bytes))
	os.Exit(1)
}

// checkError 错误检测
func checkError(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "Fatal error:%s\n", e.Error())
		os.Exit(1)
	}
}
