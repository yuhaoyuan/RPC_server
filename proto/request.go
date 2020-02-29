package proto

import "github.com/yuhaoyuan/RPC_server/dal"

var LoginRequest func(string, string) (dal.UserInfo, error)
var RegisterRequest func(userName string, pwd string) (dal.UserInfo, error)
var ModifyInfoRequest func(userName, pwd, nickName, picture string) (dal.UserInfo, error)
