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

	args := &service.Args{A: 17, B: 8}
	quotient := new(service.Quotient)
	call := client.Go("Ardith.Divide", args, &quotient, nil)
	<-call.Done

	fmt.Printf("ardith: %d / %d = %d, %d %% %d = %d\n", args.A, args.B, quotient.Quo, args.A, args.B, quotient.Rem)
}
