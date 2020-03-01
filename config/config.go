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
	AesTokenKey string
}

func BaseConfInit(){
	BaseConf.Addr = os.Getenv("ADDR")
	BaseConf.RedisAddr = os.Getenv("RedisAddr")
	BaseConf.DbAddr = os.Getenv("DbAddr")
	BaseConf.DbUser = os.Getenv("DbUser")
	BaseConf.DbPwd = os.Getenv("DbPwd")
	BaseConf.DbDatabase = os.Getenv("DbDatabase")
	BaseConf.AesTokenKey = os.Getenv("AESTOKENKEY")
}