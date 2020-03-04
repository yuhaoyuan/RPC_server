package main

import (
	"encoding/gob"
	"github.com/yuhaoyuan/RPC_server/config"
	"github.com/yuhaoyuan/RPC_server/corn"
	"github.com/yuhaoyuan/RPC_server/dal"
	"log"
	"net"
)

func init() {
	config.BaseConfInit()
}

func main() {
	gob.Register(dal.UserInfo{})
	conn, err := net.Dial("tcp", config.BaseConf.Addr)
	if err != nil {
		log.Printf("client-dial failed!")
	}
	cli := corn.NewClient(conn)
	//defer conn.Close()

	var loginRequest func(string, string) (dal.UserInfo, error)
	cli.Call("userLogin", &loginRequest)
	rsp, err := loginRequest("u0001", "12345u") // 发送请求
	if err != nil {
		log.Println(err)
	} else {
		log.Println(rsp)
	}
}
