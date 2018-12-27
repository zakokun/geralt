package main

import (
	"context"
	"log"

	"time"

	"fmt"

	"github.com/olivere/elastic"
)

var (
	esClient *elastic.Client
)

const (
	BaseType = "base"
	CppName  = "cpp"
	PHPName  = "php"
	Mapping  = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"base":{
			"properties":{
				"@timestamp":{
					"type":"date"
				}
			}
		}
	}
}`
)

func startES() (err error) {
	var (
		ctx  = context.Background()
		addr = "http://127.0.0.1:9200"
	)
	esClient, err = elastic.NewSimpleClient(elastic.SetURL(addr))
	if err != nil {
		log.Printf("NewClient(%s) err(%v)", addr, err)
		return err
	}
	_, _, err = esClient.Ping(addr).Do(ctx)
	if err != nil {
		log.Printf("Ping(%s) err(%v)", addr, err)
	}
	return
}

func listenChan() {
	var (
		err       error
		ok        bool
		app       interface{}
		indexName string
		ctx       = context.Background()
	)
	for i := range channel {
		if app, ok = i["app"]; !ok {
			log.Printf("channel get illegal data %v", i)
			continue
		}
		indexName, err = setIndexName(app.(string))
		if err != nil {
			log.Printf("setIndexName(%v) err(%v)", app, err)
			continue
		}
		ok, err := esClient.IndexExists(indexName).Do(ctx)
		if err != nil {
			log.Printf("IndexExists(%s) check err(%v)", indexName, err)
			continue
		}
		if !ok {
			if err = setIndex(indexName); err != nil {
				log.Printf("setIndex(%s) err (%v)", indexName, err)
				continue
			}
			log.Printf("create index(%s) success!", indexName)
		}
		i["@timestamp"] = time.Now().Unix() * 1000
		idxRet, err := sendLog(ctx, indexName, i)
		if err != nil {
			log.Printf("send to es err(%v)", err)
		}
		log.Printf("send to es ret(%v)", idxRet)
	}
}

// setIndexName 根据app_name 获取索引名称
func setIndexName(app string) (name string, err error) {
	var (
		day = time.Now().Format("2006-01-02")
	)
	if app == CppName || app == PHPName {
		name = app + "-" + day
	} else {
		err = fmt.Errorf("app nsetIndexame not register yet")
	}
	return
}

// setIndex 创建索引
func setIndex(name string) (err error) {
	_, err = esClient.CreateIndex(name).BodyString(Mapping).Do(context.Background())
	return
}

func sendLog(ctx context.Context, idx string, d map[string]interface{}) (idxRet *elastic.IndexResponse, err error) {
	d["@timestamp"] = time.Now().Unix() * 1000
	return esClient.Index().Index(idx).Type(BaseType).BodyJson(d).Do(ctx)
}
