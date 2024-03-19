package Cache

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"time"
)

const defaultErrorsBroadcastBuffer = 200

const redisProvider = "redis"

const defaultReconnectionTimeout = 5

type redisError struct {
	err     error
	command command
}

// GetCommand retrieves cache command associated with error as string
func (r *redisError) GetCommand() string {
	return string(r.command)
}

func (r *redisError) Error() string {
	return fmt.Sprintf("redis: error %s occured on attempt to perform command %s", r.err, r.command)
}

// RedisCache wraps *redis.Client to meet swappable and mockable Cache interface
type RedisCache struct {
	client       *redis.Client
	cacheErrors  chan Error
	broadcasting bool
	decoder      godecoder.Decoder
}

func newRedisCache(address, password string, decoder godecoder.Decoder, logger *zap.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:               address,
		Password:           password,
		IdleTimeout:        55 * time.Second,
		IdleCheckFrequency: 170 * time.Second,
	})
	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err == nil {
		redisCache := &RedisCache{client: client, decoder: decoder}
		return redisCache, nil
	}

	logger.Error("error when starting redis server", zap.Error(err))

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second * time.Duration(defaultReconnectionTimeout))

	for {
		select {

		case <-timeoutExceeded:
			return nil, fmt.Errorf("cache connection failed after %d timeout", defaultReconnectionTimeout)

		case <-ticker.C:
			err := client.Ping(ctx).Err()
			if err == nil {
				return &RedisCache{client: client, decoder: decoder}, nil
			}
			logger.Error("error when starting redis server", zap.Error(err))
		}
	}
}

func newRedisCacheWithErrorsBroadcasting(address, password string, decoder godecoder.Decoder, logger *zap.Logger) (*RedisCache, chan Error, error) {

	redisCache, err := newRedisCache(address, password, decoder, logger)
	if err != nil {
		return nil, nil, err
	}
	cacheErrors := make(chan Error, defaultErrorsBroadcastBuffer)
	redisCache.cacheErrors = cacheErrors
	redisCache.broadcasting = true

	return redisCache, redisCache.cacheErrors, nil
}

// Get retrieves value from Redis and serializes to pointer value
func (c *RedisCache) Get(ctx context.Context, key string, ptrValue interface{}) error {
	b, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		if c.broadcasting {
			go c.broadcastError(&redisError{err, get})
		}
		return err
	}
	buffer := bytes.NewBuffer(b)

	return c.decoder.Decode(buffer, ptrValue)
}

// Set takes key and value as input and setting Redis cache with this value
// shares errors to channel in case broadcasting mode is enabled
// it prevents manual error handling which seems noisy and unnecessary in client code
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expires time.Duration) {
	var data []byte
	var ok bool
	if data, ok = value.([]byte); !ok {
		var b bytes.Buffer
		err := c.decoder.Encode(&b, value)
		if err != nil && c.broadcasting {
			c.cacheErrors <- &redisError{err, set}
			return
		}
		data = b.Bytes()
	}
	if err := c.client.Set(ctx, key, data, expires).Err(); err != nil {
		if c.broadcasting {
			c.cacheErrors <- &redisError{err, set}
		}
	}
}

// Expire updates TTL for specified key
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := c.client.Expire(ctx, key, expiration).Err(); err != nil {
		if c.broadcasting {
			go c.broadcastError(&redisError{err, expire})
		}
		return err
	}

	return nil
}

// KeyExists checks whether specified key exists; returns true and nil if key is actually exists
func (c *RedisCache) KeyExists(ctx context.Context, key string) (bool, error) {

	err := c.client.Get(ctx, key).Err()
	if err == redis.Nil {
		return false, ErrCacheMiss
	}
	if err != nil && c.broadcasting {
		go c.broadcastError(&redisError{err, keyExists})
	}
	return true, err
}

func (c *RedisCache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func (c *RedisCache) broadcastError(err Error) {
	c.cacheErrors <- err
}
