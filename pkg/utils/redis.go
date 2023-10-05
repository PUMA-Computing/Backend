package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var Rdb *redis.Client

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("REDIS_PASS")

	log.Printf("Redis URL: %s", redisURL)
	log.Printf("Redis Password: %s", redisPassword)

	options := &redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       0,
	}

	Rdb = redis.NewClient(options)
	if err := Rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func IsTokenRevoked(tokenString string) (bool, error) {
	ctx := context.Background()
	exists, err := Rdb.SIsMember(ctx, "revoked_tokens", tokenString).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func RevokeToken(tokenString string) error {
	ctx := context.Background()
	_, err := Rdb.SAdd(ctx, "revoked_tokens", tokenString).Result()
	if err != nil {
		return err
	}

	return nil
}

func CloseRedis() {
	err := Rdb.Close()
	if err != nil {
		return
	}
}
