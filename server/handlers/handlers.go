package handlers

import (
	"project/internal/user"
	"project/pkg/bcrypt"
	"project/pkg/otp"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	user      *user.UserService
	telegram  Telegram
	otp       *otp.OTP
	bcrypt    *bcrypt.Bcrypt
	validator *validator.Validate
}

func NewHandler(
	user *user.UserService,
	telegram Telegram,
	otp *otp.OTP,
	bcrypt *bcrypt.Bcrypt,
) *Handler {
	return &Handler{
		user:      user,
		telegram:  telegram,
		otp:       otp,
		bcrypt:    bcrypt,
		validator: validator.New(),
	}
}
