package jwt

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"project/config"
	"project/settings"
	"strings"
	"time"

	"project/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	access  []byte
	refresh []byte

	cache *cache.Cache[string]
}

func (j *JWT) GenerateTokens(ctx context.Context, userID string) (accessToken string, refreshToken string, err error) {
	claims := make(jwt.MapClaims)
	claims["user"] = userID

	accessToken, err = j.generateToken(ctx, claims, j.access, settings.AccessTokenExpiration)
	if err != nil {
		err = fmt.Errorf("Failed to generate access token: %w", err)
		return
	}

	refreshToken, err = j.generateToken(ctx, claims, j.refresh, settings.RefreshTokenExpiration)
	if err != nil {
		err = fmt.Errorf("Failed to generate refresh token: %w", err)
		return
	}

	j.cache.Set(ctx, userID, refreshToken, settings.RefreshTokenExpiration)
	return
}

func (j *JWT) ValidateToken(ctx context.Context, tokenString string, accessKey bool) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("Invalid token signing method")
		}

		if accessKey {
			return j.access, nil
		}

		return j.refresh, nil
	})
	if err != nil {
		return "", fmt.Errorf("Failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Failed to get claims")
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return "", fmt.Errorf("Failed to get token expiration time: %w", err)
	} else if exp.Before(time.Now()) {
		return "", ErrExpired
	}

	userID, ok := claims["user"].(string)
	if !ok || userID == "" {
		return "", ErrUserGet
	}

	return userID, nil
}

func (j *JWT) RefreshToken(ctx context.Context, tokenString string) (string, string, error) {
	userID, err := j.ValidateToken(ctx, tokenString, false)
	if err != nil {
		if errors.Is(err, ErrExpired) || errors.Is(err, ErrUserGet) {
			return "", "", err
		}

		return "", "", fmt.Errorf("Failed to validate refresh token: %w", err)
	}

	refreshToken, found := j.cache.Get(ctx, userID)
	if !found {
		return "", "", errors.New("Failed to get refresh token from cache")
	}

	if !strings.EqualFold(tokenString, refreshToken) {
		return "", "", ErrInvalidRefreshToken
	}

	return j.GenerateTokens(ctx, userID)
}

func (j *JWT) generateToken(ctx context.Context, claims jwt.MapClaims, secret []byte, expireDuration time.Duration) (string, error) {
	now := time.Now()
	token := jwt.New(jwt.SigningMethodEdDSA)

	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(10 * time.Minute).Unix()

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("Error while generating token string: %w", err)
	}

	return tokenString, nil
}

func NewJWT(cfg *config.Config, bcrypt Bcrypt) (*JWT, error) {
	accessHash, err := bcrypt.Encode(context.Background(), base64.StdEncoding.EncodeToString([]byte(cfg.JWTAccessSecret)))
	if err != nil {
		return nil, fmt.Errorf("Failed to encode jwt access secret")
	}

	refreshHash, err := bcrypt.Encode(context.Background(), base64.StdEncoding.EncodeToString([]byte(cfg.JWTRefreshSecret)))
	if err != nil {
		return nil, fmt.Errorf("Failed to encode jwt refresh secret")
	}

	return &JWT{
		access:  []byte(accessHash),
		refresh: []byte(refreshHash),
		cache:   cache.NewCache[string](),
	}, nil
}
