package corn

import (
	"fmt"
	"github.com/yuhaoyuan/RPC_server/config"
	"github.com/yuhaoyuan/RPC_server/dal"
	"testing"
)


func Init(){
	//config.BaseConfInit()
	config.BaseConf.Addr = "127.0.0.1:8009"
	config.BaseConf.RedisAddr = "127.0.0.1:6379"
	config.BaseConf.DbAddr = "127.0.0.1:3306"
	config.BaseConf.DbUser = "root"
	config.BaseConf.DbPwd = "12345678"
	config.BaseConf.DbDatabase = "yhy"

	dal.CacherInit()
	dal.DbInit()
}

func TestUserRegister(t *testing.T) {
	Init()
	info, err := UserRegister("user_test1", "testpwd")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

func TestUserLogin(t *testing.T) {
	Init()
	info, err := UserLogin("user_test1", "testpwd")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)

	// 密码错误
	info2, err := UserLogin("user_test1", "testpwda")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info2)
}

func TestUserModifyInfo(t *testing.T) {
	Init()
	info, err := UserModifyInfo("user_test1", "", "阿斯顿🌹🌹", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}
