package main

import (
	"code-snippet/code/011/rpc_protocol/service"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	exit := make(chan string)
	ardith := new(service.Ardith)
	err := rpc.Register(ardith)
	if err != nil {
		log.Fatal(err.Error())
	}

	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err.Error())
	}

	go http.Serve(listener, nil)

	<-exit
}
