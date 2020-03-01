package main

import (
	"encoding/gob"
	"github.com/yuhaoyuan/RPC_server/config"
	"github.com/yuhaoyuan/RPC_server/corn"
	"github.com/yuhaoyuan/RPC_server/dal"
)

func init() {
	config.BaseConfInit()
	dal.CacherInit()
	dal.DbInit()
}

func main() {
	gob.Register(dal.UserInfo{})

	srv := corn.NewServer(config.BaseConf.Addr)
	srv.Register("userLogin", corn.UserLogin)
	srv.Register("userRegister", corn.UserRegister)
	srv.Register("UserModifyInfo", corn.UserModifyInfo)
	srv.Register("CheckToken", corn.GetUserInfoByToken)
	go srv.Run()

	for i:=0; i<10; i++{
		i--
	}
}
