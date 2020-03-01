package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/yuhaoyuan/RPC_server/config"
	"math/rand"
	"testing"
	"time"
)

func TestMakeToken(t *testing.T){
	config.BaseConfInit()
	key := []byte("wem0Upqsl5MBD0Z3")

	rand.Seed(time.Now().Unix())
	randomNumber := rand.Int63()
	sss := fmt.Sprintf("yyyhyserver%s%d%d","usertest", time.Now().Unix(), randomNumber)
	result, err := AesEncrypt([]byte(sss), []byte(config.BaseConf.AesTokenKey))
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	origData, err := AesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}