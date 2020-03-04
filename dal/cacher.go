package dal

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/yuhaoyuan/RPC_server/config"
	"log"
	"time"
)

// RedisDb redis连接
var RedisDb *redis.Client

const (
	userInfoKey = "user_info_key_%s"
)

// CacherInit 初始化redis
func CacherInit() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.BaseConf.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// checkCacherDB redis连接检查
func checkCacherDB(Db *redis.Client) error {
	err := Db.Ping().Err()
	if err != nil {
		log.Println("Error!!! Cacher DB Ping Error! info=")
		log.Println(Db.Options())
		return err
	}
	return nil
}

// CacherGetUserInfo get user info
func CacherGetUserInfo(userName string, db *redis.Client) (UserInfo, error) {
	if checkErr := checkCacherDB(db); checkErr != nil {
		return UserInfo{}, checkErr
	}
	rKey := fmt.Sprintf(userInfoKey, userName)
	userInfoString := db.Get(rKey).Val()
	if userInfoString == "" {
		return UserInfo{}, nil
	}
	var rs UserInfo
	err := json.Unmarshal([]byte(userInfoString), &rs)
	if err != nil {
		log.Printf(fmt.Sprintf("CacherGetUserInfo-Unmarshal failed!, err = %s", err))
		return rs, err
	}
	return rs, nil
}

// CacherSetUserInfo ser user info
func CacherSetUserInfo(userinfo UserInfo, db *redis.Client) error {
	if checkErr := checkCacherDB(db); checkErr != nil {
		return checkErr
	}
	rKey := fmt.Sprintf(userInfoKey, userinfo.Name)
	userInfoString, err := json.Marshal(userinfo)
	if err != nil {
		log.Printf(fmt.Sprintf("CacherSetUserInfo-Marshal failed!, err = %s", err))
		return err
	}
	err = db.Set(rKey, userInfoString, time.Minute*30).Err()
	if err != nil {
		log.Printf(fmt.Sprintf("CacherSetUserInfo-Set failed!, err = %s", err))
		return err
	}
	return nil
}

// CacherDelUserInfo del user info
func CacherDelUserInfo(userName string, db *redis.Client) error {
	if checkErr := checkCacherDB(db); checkErr != nil {
		return checkErr
	}
	rKey := fmt.Sprintf(userInfoKey, userName)
	err := db.Del(rKey).Err()
	if err != nil {
		log.Printf(fmt.Sprintf("CacherDelUserInfo-Set failed!, err = %s", err))
		return err
	}
	return nil
}
