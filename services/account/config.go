package account

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DATABASE_URL       string `envconfig:"DATABASE_URL"`
	JWT_SECRET         string `envconfig:"JWT_SECRET"`
	REFRESH_JWT_SECRET string `envconfig:"REFRESH_JWT_SECRET"`
	REDIS_URL          string `envconfig:"REDIS_URL"`
}

var AppConfig Config

func LoadConfig() {
	err := envconfig.Process("", &AppConfig)
	if err != nil {
		log.Fatal("Failed to load environment config:", err)
	}
}

func InitRedis(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return client
}
