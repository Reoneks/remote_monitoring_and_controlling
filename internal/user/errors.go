package user

import "errors"

var (
	ErrInvalidPhone    = errors.New("Invalid phone number")
	ErrInvalidPassword = errors.New("Invalid password")
	ErrInvalidOtpCode  = errors.New("Invalid otp code")
)
