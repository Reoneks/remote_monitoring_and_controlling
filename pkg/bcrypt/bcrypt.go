package bcrypt

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func (j *Bcrypt) Encode(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (j *Bcrypt) Validate(ctx context.Context, hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewBcrypt() *Bcrypt {
	return new(Bcrypt)
}
