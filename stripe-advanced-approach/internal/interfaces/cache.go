package interfaces

import "context"

// client for interfacing with the implemented cache
type CacheClient interface {
	Connect(ctx context.Context) error
	Close() error
	Ping(ctx context.Context) error
}
