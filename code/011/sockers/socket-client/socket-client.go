package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}

	service := os.Args[1]

	ipAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, ipAddr)
	checkError(err)

	conn.SetKeepAlive(true)

	_, err = conn.Write([]byte("timestamp"))
	checkError(err)

	bytes, err := ioutil.ReadAll(conn)
	checkError(err)

	fmt.Println(string(bytes))
	os.Exit(0)
}

func checkError(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", e.Error())
		os.Exit(1)
	}
}
