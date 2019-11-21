package main

import (
	"code-snippet/code/011/rpc_protocol/service"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal(err.Error())
	}

	args := &service.Args{A: 7, B: 8}
	var reply int
	err = client.Call("Ardith.Multiply", args, &reply)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("ardith: %d * %d = %d\n", args.A, args.B, reply)
}
