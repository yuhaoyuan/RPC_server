package dal

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
	"time"
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
		fmt.Println(err)
	}

	err = redisDb.Set("ttttestcacher", 1, time.Minute).Err()
	if err != nil {
		fmt.Println(err)
	}
}