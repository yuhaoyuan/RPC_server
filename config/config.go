package config

import (
	"log"
)

// BaseConf 参数
var BaseConf = BaseConfig{}

// BaseConfig config 结构体
type BaseConfig struct {
	Addr        string
	RedisAddr   string
	DbAddr      string
	DbUser      string
	DbPwd       string
	DbDatabase  string
	AesTokenKey string
}

// BaseConfInit 初始化环境变量
func BaseConfInit() {
	//BaseConf.Addr = os.Getenv("ADDR")
	//BaseConf.RedisAddr = os.Getenv("RedisAddr")
	//BaseConf.DbAddr = os.Getenv("DbAddr")
	//BaseConf.DbUser = os.Getenv("DbUser")
	//BaseConf.DbPwd = os.Getenv("DbPwd")
	//BaseConf.DbDatabase = os.Getenv("DbDatabase")
	//BaseConf.AesTokenKey = os.Getenv("AESTOKENKEY")
	BaseConf.Addr = "127.0.0.1:8009"
	BaseConf.RedisAddr = "127.0.0.1:6379"
	BaseConf.DbAddr = "127.0.0.1:3306"
	BaseConf.DbDatabase = "yhy"
	BaseConf.DbUser = "root"
	BaseConf.DbPwd = "12345678"
	BaseConf.AesTokenKey = "wem0Upqsl5MBD0Z3"

	log.Println("BaseConf = ", BaseConf)
}
