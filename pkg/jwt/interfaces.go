package jwt

import "context"

type Bcrypt interface {
	Encode(ctx context.Context, password string) (string, error)
}
