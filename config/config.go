package config

import (
	"os"
)

var BaseConf = BaseConfig{}

type BaseConfig struct {
	Addr string
	RedisAddr string
	DbAddr string
	DbUser string
	DbPwd string
	DbDatabase string
}

func BaseConfInit(){
	BaseConf.Addr = os.Getenv("ADDR")
	BaseConf.RedisAddr = os.Getenv("RedisAddr")
	BaseConf.DbAddr = os.Getenv("DbAddr")
	BaseConf.DbAddr = os.Getenv("DbUser")
	BaseConf.DbAddr = os.Getenv("DbPwd")
	BaseConf.DbDatabase = os.Getenv("DbDatabase")
}