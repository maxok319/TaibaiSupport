package TaibaiDBHelper

import (
	"github.com/go-redis/redis"
)

var redisClientInstance *redis.Client
func GetRedisClient() *redis.Client{
	return redisClientInstance
}

func init(){
	redisClientInstance = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
	})
}