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
	Mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"%s":{
			"properties":{
				"%s":{
					"type":"date",
            		"format": "strict_date_optional_time||epoch_millis"
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
		err                error
		ok                 bool
		app                interface{}
		indexName, appName string
		ctx                = context.Background()
	)
	for i := range channel {
		if app, ok = i["app"]; !ok {
			log.Printf("channel get illegal data %v", i)
			continue
		}
		appName, ok = app.(string)
		if !ok {
			log.Printf("app field must be string  data:%v", app)
			continue
		}
		indexName, err = setIndexName(appName)
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
	for _, a := range conf.Apps {
		if app == a {
			name = app + "-" + day
			return
		}
	}
	err = fmt.Errorf("app:%s not register yet", app)
	return
}

// setIndex 创建索引
func setIndex(name string) (err error) {
	_, err = esClient.CreateIndex(name).BodyString(fmt.Sprintf(Mapping, conf.TypeField, conf.TimeField)).Do(context.Background())
	return
}

func sendLog(ctx context.Context, idx string, d map[string]interface{}) (idxRet *elastic.IndexResponse, err error) {
	d["@timestamp"] = time.Now().Unix() * 1000
	return esClient.Index().Index(idx).Type(conf.TypeField).BodyJson(d).Do(ctx)
}
