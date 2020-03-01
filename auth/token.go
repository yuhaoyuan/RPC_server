package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/yuhaoyuan/RPC_server/config"
	"log"
	"math/rand"
	"time"
)

const (
	userTokenKey     = "user_token_key_%s"
	userLastTokenKey = "user_last_token_key_%s"
)

// 判断userName 和 token 是否一一对应
func CacherJudUserToken(token string, db *redis.Client) string {
	rKey := fmt.Sprintf(userTokenKey, token)
	userName := db.Get(rKey).Val()
	return userName
}

func CacherGetToken(userName string, db *redis.Client) string {
	// 先检查是否存有半小时内最近一次的token
	lastKey := fmt.Sprintf(userLastTokenKey, userName)
	token := db.Get(lastKey).Val()
	if token != "" {
		return token
	}
	// 生成token
	rand.Seed(time.Now().Unix())
	randomNumber := rand.Int63()
	orginData := fmt.Sprintf("yyyhy_server_%s_%d_%d", userName, time.Now().Unix(), randomNumber)
	cryptedToken, _ := AesEncrypt([]byte(orginData), []byte(config.BaseConf.AesToken))


	rKey := fmt.Sprintf(userTokenKey, string(cryptedToken))
	err := db.Set(rKey, userName, time.Minute*30).Err()
	if err != nil {
		log.Println("SetUserToken failed, err = ", err)
		return ""
	}
	return string(cryptedToken)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
