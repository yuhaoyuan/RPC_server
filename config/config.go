package config

import "os"

var BaseConf = BaseConfig{}

type BaseConfig struct {
	Addr string
}

func BaseConfInit(){
	BaseConf.Addr = os.Getenv("ADDR")
}