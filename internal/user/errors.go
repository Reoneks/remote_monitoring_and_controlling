package user

import "errors"

var (
	ErrInvalidPassword = errors.New("Invalid password")
	ErrInvalidOtpCode  = errors.New("Invalid otp code")
)
