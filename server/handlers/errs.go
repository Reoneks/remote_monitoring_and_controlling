package handlers

import "errors"

var (
	ErrBind     = errors.New("Invalid data")
	ErrLogin    = errors.New("Failed to login")
	ErrRegister = errors.New("Failed to register user")
)

func newHTTPError(err error) map[string]any {
	return map[string]any{
		"error": err.Error(),
	}
}
