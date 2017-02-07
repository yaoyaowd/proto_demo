package main

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"../gen-go/superroot"
)

const (
    addr = "localhost:8999"
)

func handleClient(client *superroot.SuperRootClient) (err error) {
	for i:=0; i < 1000; i++ {
		query := "shoes"
		request := superroot.SearchRequest{
			Query:&query,
			Offset:0,
			Limit:25,
		}
		client.Search(&request)
	}
	return nil
}

func main() {
	var transport thrift.TTransport
	var err error
	transport, err = thrift.NewTSocket(addr)
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport = transportFactory.GetTransport(transport)
	defer transport.Close()
	if err = transport.Open(); err != nil {
		panic(err)
	}
	handleClient(superroot.NewSuperRootClientFactory(
		transport, thrift.NewTBinaryProtocolFactoryDefault()))
}