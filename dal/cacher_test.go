package dal

import (
	"github.com/go-redis/redis"
	"testing"
)

func TestCacher(t *testing.T) {
	var redisDb *redis.Client
	redisDb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisDb.Ping().Result()
	if err != nil {
		t.Errorf("err = %v want nil", err)
	}

	userInfo := UserInfo{
		Name: "testaabbcc",
		Pwd:  "testpwd",
	}

	err = CacherSetUserInfo(userInfo, redisDb)
	if err != nil {
		t.Errorf("err = %v want nil", err)
	}
	rsp, err := CacherGetUserInfo(userInfo.Name, redisDb)
	if err != nil {
		t.Errorf("err = %v want nil", err)
	}
	if rsp.Name != userInfo.Name {
		t.Errorf("rsp.Name should = userInfo.Name")
	}
}
