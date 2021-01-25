package database

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

// InitRedis .
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping().Result()

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(pong)

}

// CreateRedisConn .
func CreateRedisConn() *redis.Client {
	return rdb
}
