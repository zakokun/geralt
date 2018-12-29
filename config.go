package main

import (
	"project/go_sdk/log"

	"github.com/BurntSushi/toml"
)

var (
	conf *Config
)

type Config struct {
	Apps      []string
	TimeField string
	TypeField string
}

func initConfig() (err error) {
	conf = new(Config)
	_, err = toml.DecodeFile("./config.toml", conf)
	if err != nil {
		log.Error("toml.DecodeFile(%s) err(%v)", "./config.toml", err)
	}
	if conf.TimeField == "" {
		conf.TimeField = "@timestamp"
	}
	if conf.TypeField == "" {
		conf.TypeField = "base"
	}
	return
}
