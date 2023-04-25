package jwt

import "errors"

var (
	ErrExpired = errors.New("Token already expired")
	ErrUserGet = errors.New("Failed to get user id from token")
)
