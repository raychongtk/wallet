package datastore

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/raychongtk/wallet/util"
	"log"
)

func ProvideRedis() redis.Client {
	return *redisConnection()
}

func redisConnection() *redis.Client {
	var config *RedisConfig
	util.Load("dev", ".", &config)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Test connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")

	return redisClient
}

type RedisConfig struct {
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}
