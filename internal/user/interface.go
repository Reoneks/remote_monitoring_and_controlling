package user

import (
	"context"
	"remote_monitoring_and_controlling/pkg/postgres"
)

type DB interface {
	GetUserByPhone(ctx context.Context, phone string) (postgres.User, error)
	CreateUser(ctx context.Context, user *postgres.User) error
	AddContactInfo(ctx context.Context, contactInfo *postgres.ContactInfo) error
	EnableOTP(ctx context.Context, userID, secret string) error
	DisableOTP(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, userID string) error
}
