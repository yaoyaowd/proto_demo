package main

import (
    "fmt"
    "git.apache.org/thrift.git/lib/go/thrift"

    "./handler"
    "../gen-go/superroot"
)

const (
    addr = "localhost:8999"
)

func main() {
    var transport thrift.TServerTransport
    var err error
    transport, err = thrift.NewTServerSocket(addr)
    if err != nil {
        panic(err)
    }

    handler := handler.NewSuperRootHandler()
    processor := superroot.NewSuperRootProcessor(handler)
    transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

    server := thrift.NewTSimpleServer4(
        processor,
        transport,
        transportFactory,
        protocolFactory)

    fmt.Println("Starting the simple server... on ", addr)
    server.Serve()
}