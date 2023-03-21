package settings

import "time"

const (
	AccessTokenExpiration  = 10 * time.Minute
	RefreshTokenExpiration = 20 * time.Minute
)
