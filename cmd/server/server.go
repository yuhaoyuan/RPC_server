package main

import (
	"encoding/gob"
	"github.com/yuhaoyuan/RPC_server/config"
	"github.com/yuhaoyuan/RPC_server/corn"
	"github.com/yuhaoyuan/RPC_server/dal"
	"github.com/yuhaoyuan/RPC_server/proto"
	"time"
)

func init() {
	config.BaseConfInit()
	dal.CacherInit()
	dal.DbInit()
}

func main() {
	gob.Register(proto.User{})

	srv := corn.NewServer(config.BaseConf.Addr)
	srv.Register("userLogin", corn.UserLogin)
	srv.Register("userRegister", corn.UserRegister)
	srv.Register("UserModifyInfo", corn.UserModifyInfo)
	go srv.Run()
	time.Sleep(time.Minute * 5)
}
