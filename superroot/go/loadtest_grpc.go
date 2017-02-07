package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync/atomic"
	pb "./superroot"
	"time"
)

const (
	total_requests = 100000
	test_qps = 2000
	address = "localhost:8999"
)

func loadtest(c *pb.SuperRootClient) (error) {
	var total_latency, processed_requests, total_received_docs int64
	total_latency = 0
	processed_requests = 0
	total_received_docs = 0
	wait := int64(float64(time.Second) / test_qps)

	for i := 0; i < total_requests; i++ {
		time.Sleep(time.Duration(wait))
		go func() {
			start_time := int64(time.Now().UnixNano() / 1000000)
			ret, err := (*c).Search(
				context.Background(),
				&pb.SearchRequest{
					Query:"shoes",
					Offset:0,
					Limit:25})
			if err == nil {
				atomic.AddInt64(&total_received_docs, int64(len(ret.Docs)))
			}
			end_time := int64(time.Now().UnixNano() / 1000000)
			time_spend := end_time - start_time

			tl := atomic.AddInt64(&total_latency, time_spend)
			cr := atomic.AddInt64(&processed_requests, 1)
			if cr % 1000 == 0 {
				log.Printf("Send %d requests\n", i)
				log.Printf("Processed %d requests % docs\n", cr, total_received_docs)
				log.Printf("Avg latency: %fms\n", float64(tl) / float64(cr))
			}
		}()
	}
	return nil
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSuperRootClient(conn)

	err = loadtest(&c)
	if err != nil {
		log.Fatal("could not search: %v", err)
	}
}
