package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"sort"
	pb "./superroot"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8999"
)

// Put your solr service here.
var hosts = []string{}

type Server struct{}

// Search a single solr host.
func search(host string, query string, offset int32, limit int32, ch chan<- []pb.SearchDoc) {
	// This is for demo, use &url.URL{Path: request string} to avoid url encode / decode problems.
	url_str := fmt.Sprintf("http://%s:8999/solr/wishsolrcluster/select?" +
		"defType=edismax&qf=product_description%%20text&" +
		"wt=json&start=%d&count=%d&docsAllowed=50000&omitHeader=true" +
		"&fl=id%%20score&q={!cache=false}%s", host, offset, limit, query)
	url, err := url.Parse(url_str)
	if err != nil {
		ch <- make([]pb.SearchDoc, 0)
	}

	resp, err := http.Get(url.String())
	if err != nil {
		ch <- make([]pb.SearchDoc, 0)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result interface{}
	json.Unmarshal(body, &result)

	response := (result).(map[string]interface{})["response"]
	docs := response.(map[string]interface{})["docs"].([]interface{})
	ret := []pb.SearchDoc{}
	for i := 0; i < len(docs); i++ {
		doc := docs[i].(map[string]interface{})
		ret = append(ret, pb.SearchDoc{
			Id: doc["id"].(string),
			Score: float32(doc["score"].(float64)),
		})
	}
	ch <- ret
}

// Helper structure to sort search document by score.
type Docs []*pb.SearchDoc
func (c Docs) Len() int {
	return len(c)
}
func (c Docs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Docs) Less(i, j int) bool {
	return c[i].Score > c[j].Score
}

func (s *Server) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	ch := make(chan []pb.SearchDoc)
	for _, host := range hosts {
		go search(host, in.Query, 0, in.Offset+in.Limit, ch)
	}

	ret := []*pb.SearchDoc{}
	for range hosts {
		host_ret := <- ch
		for _, doc := range host_ret {
			ret = append(ret, &doc)
		}
	}

	sort.Sort(Docs(ret))
	return &pb.SearchResponse{Docs: ret[in.Offset:in.Offset + in.Limit]}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSuperRootServer(s, &Server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
