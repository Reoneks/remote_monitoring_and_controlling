package otp

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"remote_monitoring_and_controlling/settings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type OTP struct{}

func (t *OTP) GenerateKey(ctx context.Context, userID string) ([]byte, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "A remote system for monitoring and controlling business processes and tasks",
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

func (t *OTP) ValidateKey(ctx context.Context, key, secret string) bool {
	return totp.Validate(key, secret)
}

func NewOTP() *OTP {
	return new(OTP)
}
