package handlers

import "context"

type Postgres interface{}

type Telegram interface {
	SendNotification(ctx context.Context, text string, userID int64) error
}

type OTP interface {
	GenerateKey(ctx context.Context, userID string) ([]byte, string, error)
	ValidateKey(ctx context.Context, key, secret string) (bool, error)
}

type Bcrypt interface {
	Encode(ctx context.Context, password string) (string, error)
	Validate(ctx context.Context, hash, password string) bool
}
