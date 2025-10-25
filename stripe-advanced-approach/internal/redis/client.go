package redis

import (
	"context"
	"encoding/json"
	"fmt"
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
	// Test the connection
	_, err := c.rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Printf("Redis client connected successfully to %s\n", util.GetEnv("REDIS_ADDR", "localhost:6379"))
	return nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	result := c.rdb.SetNX(ctx, key, value, expiration)
	return result.Val(), result.Err()
}

func (c *Client) Pipeline() redislib.Pipeliner {
	return c.rdb.Pipeline()
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result := c.rdb.Exists(ctx, key)
	return result.Val() > 0, result.Err()
}

// TTL returns the remaining time to live of a key
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.rdb.TTL(ctx, key).Result()
}

// Expire sets a timeout on a key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.rdb.Expire(ctx, key, expiration).Err()
}

// Keys finds all keys matching a pattern (use carefully in production)
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.rdb.Keys(ctx, pattern).Result()
}

// GetSetJSON is a helper for JSON operations
func (c *Client) GetSetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return c.Set(ctx, key, jsonData, expiration)
}

// GetJSON is a helper for JSON retrieval
func (c *Client) GetJSON(ctx context.Context, key string, dest interface{}) error {
	jsonStr, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonStr), dest)
}

/**
* Helper methods to assist in maintaining the synchronization of keys
*
* There are 3 primary keys for mappings:
*
* -- Data --
* - 1. stripe:customer:{customerId} → full customer stripe data
*
* -- Key Mappings --
* - 2. stripe:customer:{customerId}:userid →  userId (customerId →  userId lookup)
* - 3. stripe:customer:userid:{userId} →  customerId (userId →  customerId lookup)
**/

const (
	cacheKeyCustomerData       = "stripe:customer:%s"
	cacheKeyCustomerIdToUserId = "stripe:customer:%s:userid"
	cacheKeyUserIdToCustomerId = "stripe:customer:userid:%s"
)

func (c *Client) GetCustomerDataFromCustomerIdKey(customerId string) string {
	key := fmt.Sprintf(cacheKeyCustomerData, customerId)
	return key
}

// customer to userId
func (c *Client) GetUserIdFromCusIdKey(customerId string) string {
	key := fmt.Sprintf(cacheKeyCustomerIdToUserId, customerId)
	return key
}

// userId to customerId
func (c *Client) GetCusIdFromUserIdKey(userId string) string {
	key := fmt.Sprintf(cacheKeyUserIdToCustomerId, userId)
	return key
}
