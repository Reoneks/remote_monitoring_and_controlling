package cache

import (
	"context"
	"project/settings"
	"time"

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

	if res, ok := val.(T); !ok {
		return res, false
	} else {
		return res, found
	}
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{c: cache.New(settings.DefaultCacheExpiration, settings.CacheCleanup)}
}
