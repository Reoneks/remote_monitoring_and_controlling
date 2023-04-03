package handlers

import (
	"remote_monitoring_and_controlling/internal/user"
	"remote_monitoring_and_controlling/pkg/bcrypt"
	"remote_monitoring_and_controlling/pkg/otp"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	user      *user.UserService
	otp       *otp.OTP
	bcrypt    *bcrypt.Bcrypt
	validator *validator.Validate
}

func NewHandler(
	user *user.UserService,
	otp *otp.OTP,
	bcrypt *bcrypt.Bcrypt,
) *Handler {
	return &Handler{
		user:      user,
		otp:       otp,
		bcrypt:    bcrypt,
		validator: validator.New(),
	}
}
