package cache

import (
	"context"
	"time"

	"remote_monitoring_and_controlling/settings"

	"github.com/patrickmn/go-cache"
)

type Cache[T any] struct {
	c *cache.Cache
}

func (c *Cache[T]) Set(ctx context.Context, key string, value T, duration time.Duration) {
	c.c.Set(key, value, duration)
}

func (c *Cache[T]) Get(ctx context.Context, key string) (T, bool) {
	val, found := c.c.Get(key)
	if !found {
		var res T
		return res, found
	}

	res, ok := val.(T)
	if !ok {
		return res, false
	}

	return res, found
}

func (c *Cache[T]) Delete(ctx context.Context, key string) {
	c.c.Delete(key)
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{c: cache.New(settings.DefaultCacheExpiration, settings.CacheCleanup)}
}
