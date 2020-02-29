package dal

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/yuhaoyuan/RPC_server/config"
	"log"
	"time"
)

var RedisDb *redis.Client

const (
	userInfoKey = "user_info_key_%s"
)

func CacherInit() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.BaseConf.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func checkCacherDB(Db *redis.Client) error {
	err := Db.Ping().Err()
	if err != nil {
		log.Println("Error!!! Cacher DB Ping Error! info=")
		log.Println(Db.Options())
		return err
	}
	return nil
}

func CacherGetUserInfo(userName string, db *redis.Client) (UserInfo, error) {
	if checkErr := checkCacherDB(db); checkErr!= nil {
		return UserInfo{}, checkErr
	}
	rKey := fmt.Sprintf(userInfoKey, userName)
	userInfoString := db.Get(rKey).Val()
	if userInfoString == "" {
		return UserInfo{}, nil
	} else {
		var rs UserInfo
		err := json.Unmarshal([]byte(userInfoString), &rs)
		if err != nil {
			log.Printf(fmt.Sprintf("CacherGetUserInfo-Unmarshal failed!, err = %s", err))
			return rs, err
		}
		return rs, nil
	}
}

func CacherSetUserInfo(userinfo UserInfo, db *redis.Client) error {
	if checkErr := checkCacherDB(db); checkErr!= nil {
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
