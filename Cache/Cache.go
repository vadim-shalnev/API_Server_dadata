package Cache

import (
	"context"
	"errors"
	"time"
)

type command string

const (
	get       command = "get"
	set       command = "set"
	keyExists command = "keyExists"
	expire    command = "expire"
)

var (
	// ErrCacheMiss is a replacement for implementation defined cache miss error of different providers,
	// such as replacement for redis.Nil which looks meaningless
	ErrCacheMiss = errors.New("cache: key not found")
)

// Cache is an interface which abstracts cache providers details for mockability and usability
type Cache interface {
	// Retrieves the content associated with the given key. decoding it into the given
	// pointer.
	//
	// Returns:
	//   - nil if the value was successfully retrieved and ptrValue set
	//   - ErrCacheMiss if the value was not in the cache
	//   - an implementation specific error otherwise
	Get(ctx context.Context, key string, ptrValue interface{}) error

	// set the given key/value in the cache, overwriting any existing value
	// associated with that key
	// if broadcasting mode is specified sends related errors to returned channel during cache initialization
	Set(ctx context.Context, key string, value interface{}, expires time.Duration)

	// commented for initial commit
	Delete(key string) error
	// Add(key string, value interface{}, expires time.Duration) error // only if the key does not exist
	// Replace(key string, value interface{}, expires time.Duration) error // only if the key already exists

	// verifies whether specified key exists in cache; returns true in case key exists and nil error
	KeyExists(ctx context.Context, key string) (bool, error)

	// use this when you want time to live of a key will be updated to the new value
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// Error is an interface which helps to communicate cache errors for logging and metrics
type Error interface {
	error
	GetCommand() string
}
