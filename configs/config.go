package configs

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisURL  string
	RedisPass string

	ApiPort      string
	JWTSecretKey string

	CloudflareAccountId   string
	CloudflareR2AccessId  string
	CloudflareR2AccessKey string

	AWSAccessKeyId     string
	AWSSecretAccessKey string
	AWSRegion          string
	S3Bucket           string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		DBUser:                os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBName:                os.Getenv("DB_NAME"),
		RedisURL:              os.Getenv("REDIS_URL"),
		RedisPass:             os.Getenv("REDIS_PASS"),
		ApiPort:               os.Getenv("API_PORT"),
		JWTSecretKey:          os.Getenv("JWT_SECRET_KEY"),
		CloudflareAccountId:   os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		CloudflareR2AccessId:  os.Getenv("CLOUDFLARE_R2_ACCESS_ID"),
		CloudflareR2AccessKey: os.Getenv("CLOUDFLARE_R2_ACCESS_KEY"),
		AWSAccessKeyId:        os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:             os.Getenv("AWS_REGION"),
		S3Bucket:              os.Getenv("S3_BUCKET"),
	}

	fmt.Printf("Loaded Config: %+v\n", cfg)
	return cfg
}
