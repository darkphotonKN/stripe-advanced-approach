package redis

import (
	"context"
	"fmt"

	"github.com/darkphotonKN/stripe-advanced-approach/config"
)

type Client struct {
	config *config.RedisConfig
}

func NewClient(cfg *config.RedisConfig) *Client {
	return &Client{
		config: cfg,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	fmt.Printf("Redis client configured for %s\n", c.config.GetAddr())
	return nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.config.Ping(ctx)
}

