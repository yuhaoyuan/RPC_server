package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"time"
)

const (
	userTokenKey     = "user_token_key_%s"
	userLastTokenKey = "user_last_token_key_%s"
)

// CacherGetUserNameByToken 从redis中获得user name
func CacherGetUserNameByToken(token string, db *redis.Client) string {
	rKey := fmt.Sprintf(userTokenKey, token)
	userName := db.Get(rKey).Val()
	return userName
}

// CacherGetToken 从redis中获得token
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
	orginData := fmt.Sprintf("yyyhyserver%s%d%d", userName, time.Now().Unix(), randomNumber)

	rKey := fmt.Sprintf(userTokenKey, string(orginData))
	err := db.Set(rKey, userName, time.Minute*30).Err()
	if err != nil {
		log.Println("SetUserToken failed, err = ", err)
		return ""
	}
	db.Set(lastKey, orginData, time.Minute*30)
	return string(orginData)
}

// PKCS7Padding 加密前填充数据
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 解密前处理数据
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt 将token加密
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
	return []byte(hex.EncodeToString(crypted)), nil
}

// AesDecrypt 解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	d, _ := hex.DecodeString(string(crypted))
	origData := make([]byte, len(d))
	blockMode.CryptBlocks(origData, d)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
