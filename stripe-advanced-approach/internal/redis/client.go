package redis

import (
	"context"
	"time"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/util"
	redislib "github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redislib.Client
}

func NewClient() *Client {
	rdb := redislib.NewClient(&redislib.Options{
		Addr:            util.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password:        util.GetEnv("REDIS_PASSWORD", ""),
		DB:              util.GetEnvAsInt("REDIS_DB", 0),
		MaxRetries:      util.GetEnvAsInt("REDIS_MAX_RETRIES", 3),
		PoolSize:        util.GetEnvAsInt("REDIS_POOL_SIZE", 10),
		MinIdleConns:    util.GetEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		ConnMaxIdleTime: 300 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
	})

	return &Client{
		rdb: rdb,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	// connect to redis
	return nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}
