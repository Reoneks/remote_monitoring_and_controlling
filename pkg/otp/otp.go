package otp

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"project/settings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type OTP struct{}

func (t *OTP) GenerateKey(ctx context.Context, userID string) ([]byte, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Project",
		AccountName: userID,
		Period:      settings.OTPExpiration,
		SecretSize:  settings.SecretSize,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA512,
	})
	if err != nil {
		return nil, "", fmt.Errorf("Failed to generate totp key: %w", err)
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to convert totp key into image: %w", err)
	}

	png.Encode(&buf, img)
	return buf.Bytes(), key.Secret(), nil
}

func (t *OTP) ValidateKey(ctx context.Context, key, secret string) (bool, error) {
	return totp.ValidateCustom(key, secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    settings.OTPExpiration,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
}

func NewOTP() *OTP {
	return new(OTP)
}
