package redis

import (
	"chatapp-api/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(config *config.Config) *redis.Client {
	// 1. Build redis address from config
	address := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)

	// 2. Create redis client
	client := redis.NewClient(&redis.Options{
		Addr: address,
		Password: config.Redis.Password,
		DB: 0,
	})

	// 3. Ping redis to check connection (5 seconds timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully!")
	return client
}