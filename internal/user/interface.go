package user

import (
	"context"
	"project/structs"
)

type DB interface {
	GetUserByID(ctx context.Context, userID string) (structs.User, error)
	GetUserByPhone(ctx context.Context, phone string) (structs.User, error)
	CreateUser(ctx context.Context, user *structs.User) error
	ChangePassword(ctx context.Context, userID, password string) error
	EnableOTP(ctx context.Context, userID, secret string) error
	DisableOTP(ctx context.Context, userID string) error
}
