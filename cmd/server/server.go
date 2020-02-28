package main

import (
	"encoding/gob"
	"github.com/yuhaoyuan/RPC_server/config"
	"github.com/yuhaoyuan/RPC_server/corn"
	"github.com/yuhaoyuan/RPC_server/proto"
	"time"
)

func init() {
	config.BaseConfInit()
}

func main() {
	gob.Register(proto.User{})

	srv := corn.NewServer(config.BaseConf.Addr)
	srv.Register("userLogin", corn.UserLogin)
	go srv.Run()
	time.Sleep(time.Minute * 5)
}
