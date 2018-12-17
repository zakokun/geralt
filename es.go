package main

import (
	"context"
	"fmt"
	"log"

	"github.com/olivere/elastic"
)

var (
	esClient *elastic.Client
)

func startES() (err error) {
	var (
		code int
		ctx  context.Context
		info *elastic.PingResult
		addr = "http://127.0.0.1:9200"
	)
	esClient, err = elastic.NewClient(elastic.SetURL(addr))
	if err != nil {
		log.Println(err)
		return err
	}

	info, code, err = esClient.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("Elasticsearch started with code %d and version %s\n", code, info.Version.Number)
	go dealData()
	return
}
