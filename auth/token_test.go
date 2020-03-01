package auth

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestMakeToken(t *testing.T){
	key := []byte("wem0Upqsl5MBD0Z3")

	result, err := AesEncrypt([]byte("hello world"), key)
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