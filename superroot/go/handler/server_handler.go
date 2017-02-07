package handler

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "net/url"

    sr_thrift "../../gen-go/superroot"
)

var hosts = [5]string{
	"10.10.32.23",
	"10.10.32.89",
	"10.10.32.185",
	"10.10.32.42",
	"10.10.32.165"}

type SuperRootHandler struct {
}

func NewSuperRootHandler() *SuperRootHandler {
    return &SuperRootHandler{}
}

func search(host string, query string, offset int32, limit int32, ch chan<- []sr_thrift.SearchDoc) {
	url_str := fmt.Sprintf("http://%s:8999/solr/wishsolrcluster/select?" +
		"defType=edismax&qf=product_description%%20text&" +
		"wt=json&start=%d&count=%d&docsAllowed=50000&omitHeader=true" +
		"&fl=id%%20score&q={!cache=false}%s", host, offset, limit, query)
	url, err := url.Parse(url_str)
	if err != nil {
		ch <- make([]sr_thrift.SearchDoc, 0)
	}

	resp, err := http.Get(url.String())
	if err != nil {
		ch <- make([]sr_thrift.SearchDoc, 0)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var result interface{}
	json.Unmarshal(body, &result)

	response := (result).(map[string]interface{})["response"]
	docs := response.(map[string]interface{})["docs"].([]interface{})
	ret := []sr_thrift.SearchDoc{}
	for i := 0; i < len(docs); i++ {
		jsonDoc := docs[i].(map[string]interface{})
		id := jsonDoc["id"].(string)
		score := jsonDoc["score"].(float64)
		var doc sr_thrift.SearchDoc
		doc.ID = &id
		doc.Score = &score
		ret = append(ret, doc)
	}
	ch <- ret
}

func (p *SuperRootHandler)  Search(in *sr_thrift.SearchRequest) (sr *sr_thrift.SearchResponse, err error) {
	ch := make(chan []sr_thrift.SearchDoc)
	for _, host := range hosts {
		go search(host, *in.Query, in.Offset, in.Limit, ch)
    	}
	ret := []*sr_thrift.SearchDoc{}
	for range hosts {
		host_ret := <-ch
		for _, doc := range host_ret {
			ret = append(ret, &doc)
		}
	}
	return &sr_thrift.SearchResponse{Docs:ret}, nil
}
