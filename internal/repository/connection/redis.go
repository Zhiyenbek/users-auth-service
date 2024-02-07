package connection

import (
	"log"
	"os"
	"strconv"

	"github.com/Zhiyenbek/users-auth-service/config"
	"github.com/go-redis/redis/v7"
)

func NewRedis(cfg *config.RedisConf) (*redis.Client, error) {
	var client *redis.Client

	redis_URL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		log.Println("Couldn't get redis url. Continuing with config")
		client = redis.NewClient(&redis.Options{
			Addr: cfg.Host + ":" + strconv.Itoa(cfg.Port),
			DB:   cfg.DB,
		})
	} else {
		log.Println("Redis url: ", redis_URL)
		opt, err := redis.ParseURL(redis_URL)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		client = redis.NewClient(opt)
	}

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
