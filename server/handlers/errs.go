package handlers

import (
	"errors"

	"remote_monitoring_and_controlling/internal/user"
)

var (
	ErrBind             = errors.New("Invalid data")
	ErrLogin            = errors.New("Failed to login")
	ErrRegister         = errors.New("Failed to register user")
	ErrAddPhone         = errors.New("Failed to add alternative phone number")
	ErrGetTasks         = errors.New("Failed to get tasks")
	ErrGetVacations     = errors.New("Failed to get vacations")
	ErrGetPayments      = errors.New("Failed to get payments")
	ErrSaveTasks        = errors.New("Failed to save tasks")
	ErrSaveVacations    = errors.New("Failed to save vacations")
	ErrSavePayments     = errors.New("Failed to save payments")
	ErrUpdateTaskStatus = errors.New("Failed to update task status")
)

func newHTTPError(err error) map[string]any {
	return map[string]any{
		"error": err.Error(),
	}
}

func checkErr(err, defaultErr error) error {
	switch {
	case errors.Is(err, user.ErrInvalidOtpCode):
		return user.ErrInvalidOtpCode
	case errors.Is(err, user.ErrInvalidPassword):
		return user.ErrInvalidPassword
	case errors.Is(err, user.ErrInvalidPhone):
		return user.ErrInvalidPhone
	case errors.Is(err, user.ErrInvalid2FAID):
		return user.ErrInvalid2FAID
	default:
		return defaultErr
	}
}
