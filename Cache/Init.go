package Cache

import (
	"fmt"
	"github.com/ptflp/godecoder"
	"gitlab.com/ptflp/goboilerplate/config"
	"go.uber.org/zap"
)

// NewCache takes config and depending on cache section of config wraps actual cache implementation
func NewCache(cacheConfig config.Cache, decoder godecoder.Decoder, logger *zap.Logger) (Cache, error) {
	redisCache, err := newRedisCache(fmt.Sprintf("%s:%s", cacheConfig.Address, cacheConfig.Port), cacheConfig.Password, decoder, logger)
	if err != nil {
		return nil, err
	}

	return redisCache, nil
}

// NewCacheWithErrorsBroadcasting behaves as NewCache but additionally enables error broadcasting mode and returns channel
// caller can listen to this channel to receive errors connected with cache without explicitly manually handling errors in client code
func NewCacheWithErrorsBroadcasting(cacheConfig config.Cache, decoder godecoder.Decoder, logger *zap.Logger) (Cache, chan Error, error) {
	redisCache, cacheErrorsChannel, err := newRedisCacheWithErrorsBroadcasting(
		cacheConfig.Address, cacheConfig.Password, decoder, logger)
	if err != nil {
		return nil, nil, err
	}

	return redisCache, cacheErrorsChannel, nil
}
