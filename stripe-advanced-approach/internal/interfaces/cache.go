package interfaces

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// client for interfacing with the implemented cache
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Pipeline() redis.Pipeliner
	Close() error
	Ping(ctx context.Context) error
	GetUserIdFromCustomerIdKey(customerId string) string
	GetCustomerIdFromUserIdKey(userId string) string
	GetCustomerDataFromCustomerIdKey(customerId string) string
}
