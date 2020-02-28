package corn

import (
	"github.com/yuhaoyuan/RPC_server/proto"
	"log"
)

func UserLogin(name string, pwd string) (proto.User, error) {
	log.Printf("user-login!")
	return proto.User{
		Name: name,
		NickName: "ttt",
		Age: 18,
		Description: "18-ttt",
	}, nil
}
