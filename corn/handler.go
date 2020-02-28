package corn

import (
	"github.com/yuhaoyuan/RPC_server/proto"
	"log"
)

func UserLogin(name string, pwd string) (proto.User, error) {
	log.Printf("user-login!")

	// 先检查redis里面是否存在此用户


	return proto.User{
		Name: name,
		NickName: "ttt",
		Age: 18,
		Description: "18-ttt",
	}, nil
}
