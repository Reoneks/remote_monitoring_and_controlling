package settings

import "time"

const (
	DefaultCacheExpiration time.Duration = -1
	CacheCleanup                         = 5 * time.Minute
)
