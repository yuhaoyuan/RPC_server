package util

import (
	"fmt"
	"github.com/yuhaoyuan/RPC_server/config"
	"github.com/yuhaoyuan/RPC_server/dal"
	"log"
	"math/rand"
	"testing"
	"time"
)

var r *rand.Rand


func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := 1
		for{
			b = r.Intn(52)
			if b<=25|| b>=32{ // 27~33不是英文字母
				break
			}
		}
		b = b+65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// 看一下assic码表
func TestRandString(t *testing.T){
	for i:=25+65;i<=32+65;i++{
		bbb := byte(i)
		fmt.Println(string(bbb))
	}
}

func TestInserUserInfo(t *testing.T){
	config.BaseConfInit()
	dal.DbInit()
	r = rand.New(rand.NewSource(time.Now().Unix()))

	// 这是一个插入10,000,000条用户信息的脚本, 谨慎运行
	for i:=0; i<10000000; i++{
		fmt.Println("i = ", i)
		userName := RandString(10)
		pwd := "testpwd"
		picture := "http://q6gy4v9f7.bkt.clouddn.com/usertestc_1583059444997873000.jpg"

		inserSQL := fmt.Sprintf("insert into user_info set user_name='%s', pwd='%s', nick_name='', picture='%s'", userName, pwd, picture)
		_, err := dal.SQLDB.Exec(inserSQL)
		if err != nil {
			log.Println("DbInsertUserInfo, inserSql-error, err = ", err)
		}
	}

}