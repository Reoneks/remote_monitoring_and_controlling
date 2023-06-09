package jwt

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"remote_monitoring_and_controlling/config"
	"remote_monitoring_and_controlling/settings"

	"remote_monitoring_and_controlling/pkg/bcrypt"
	"remote_monitoring_and_controlling/pkg/cache"

	"github.com/dchest/uniuri"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sethvargo/go-password/password"
)

type JWT struct {
	secret string

	saltCache *cache.Cache[string]
}

func (j *JWT) GenerateToken(ctx context.Context, userID string) (accessToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := make(jwt.MapClaims)
	claims["user"] = userID
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	salt, err := password.Generate(10, 3, 3, false, true)
	if err != nil {
		salt = uniuri.NewLen(10)
	}

	accessToken, err = token.SignedString([]byte(j.secret + salt))
	if err != nil {
		return "", fmt.Errorf("Error while generating token string: %w", err)
	}

	j.saltCache.Set(ctx, accessToken, salt, settings.TokenDefaultExpiration)
	return
}

func (j *JWT) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid token signing method")
		}

		salt, found := j.saltCache.Get(ctx, tokenString)
		if !found {
			return nil, errors.New("Salt not found")
		}

		return []byte(j.secret + salt), nil
	})
	if err != nil {
		return "", fmt.Errorf("Failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Failed to get claims")
	}

	userID, ok := claims["user"].(string)
	if !ok || userID == "" {
		return "", ErrUserGet
	}

	return userID, nil
}

func (j *JWT) DeleteSalt(ctx context.Context, token string) {
	j.saltCache.Delete(ctx, token)
}

func NewJWT(cfg *config.Config, bcrypt *bcrypt.Bcrypt) (*JWT, error) {
	accessHash, err := bcrypt.Encode(context.Background(), base64.StdEncoding.EncodeToString([]byte(cfg.JWTSecret)))
	if err != nil {
		return nil, fmt.Errorf("Failed to encode jwt access secret")
	}

	return &JWT{secret: accessHash, saltCache: cache.NewCache[string]()}, nil
}
