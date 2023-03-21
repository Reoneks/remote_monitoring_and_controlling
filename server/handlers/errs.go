package handlers

import "errors"

var (
	ErrBind          = errors.New("Invalid data")
)

func newHTTPError(err error) map[string]any {
	return map[string]any{
		"error": err.Error(),
	}
}
