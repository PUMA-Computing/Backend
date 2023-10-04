package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
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
