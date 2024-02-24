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

	ServerPort   string
	JWTSecretKey string

	CloudflareAccountId   string
	CloudflareR2AccessId  string
	CloudflareR2AccessKey string
	CloudflareR2Bucket    string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		DBUser:                os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBName:                os.Getenv("DB_NAME"),
		ServerPort:            os.Getenv("API_PORT"),
		JWTSecretKey:          os.Getenv("JWT_SECRET_KEY"),
		CloudflareAccountId:   os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		CloudflareR2AccessId:  os.Getenv("CLOUDFLARE_R2_ACCESS_ID"),
		CloudflareR2AccessKey: os.Getenv("CLOUDFLARE_R2_ACCESS_KEY"),
		CloudflareR2Bucket:    os.Getenv("CLOUDFLARE_R2_BUCKET"),
	}

	fmt.Printf("Loaded Config: %+v\n", cfg)
	return cfg
}
