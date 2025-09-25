package config

import (
	"context"
	"fmt"
	"os"
	"time"
)

type RedisConfig struct {
	Host            string
	Port            string
	Password        string
	DB              int
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolSize        int
	MinIdleConns    int
	MaxConnAge      time.Duration
	PoolTimeout     time.Duration
	IdleTimeout     time.Duration
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:            getEnv("REDIS_HOST", "localhost"),
		Port:            getEnv("REDIS_PORT", "6379"),
		Password:        getEnv("REDIS_PASSWORD", ""),
		DB:              0,
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolSize:        10,
		MinIdleConns:    5,
		MaxConnAge:      0,
		PoolTimeout:     4 * time.Second,
		IdleTimeout:     5 * time.Minute,
	}
}

func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *RedisConfig) Ping(ctx context.Context) error {
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}