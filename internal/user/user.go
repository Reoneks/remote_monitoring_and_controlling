package user

import (
	"context"

	"remote_monitoring_and_controlling/pkg/bcrypt"
	"remote_monitoring_and_controlling/pkg/jwt"
	"remote_monitoring_and_controlling/pkg/otp"
	"remote_monitoring_and_controlling/structs"

	"github.com/dongri/phonenumber"
	"github.com/oklog/ulid/v2"
)

type Service struct {
	db     DB
	jwt    *jwt.JWT
	otp    *otp.OTP
	bcrypt *bcrypt.Bcrypt
}

func (u *Service) Register(ctx context.Context, req *structs.Register) ([]byte, error) {
	var (
		id     = ulid.Make().String()
		image  []byte
		secret string
		err    error
	)

	phoneNumber := phonenumber.Parse(req.Phone, "")
	if phoneNumber == "" {
		return nil, ErrInvalidPhone
	}

	if req.OTPEnabled {
		image, secret, err = u.otp.GenerateKey(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	encryptedPassword, err := u.bcrypt.Encode(ctx, req.Password)
	if err != nil {
		return nil, err
	}

	return image, u.db.CreateUser(ctx, &structs.User{
		ID:         id,
		Phone:      phoneNumber,
		Password:   encryptedPassword,
		OTPEnabled: req.OTPEnabled,
		OTPSecret:  secret,
	})
}

func (u *Service) Login(ctx context.Context, req *structs.Login) (string, bool, error) {
	user, err := u.db.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		return "", false, err
	}

	if !u.bcrypt.Validate(ctx, user.Password, req.Password) {
		return "", false, ErrInvalidPassword
	}

	if user.OTPEnabled {
		return "", true, nil
	}

	token, err := u.jwt.GenerateToken(ctx, user.ID)
	return token, false, err
}

func (u *Service) OTPCheck(ctx context.Context, req *structs.TwoFA) (string, error) {
	user, err := u.db.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		return "", err
	}

	ok := u.otp.ValidateKey(ctx, req.OTPPassword, user.OTPSecret)
	if !ok {
		return "", ErrInvalidOtpCode
	}

	return u.jwt.GenerateToken(ctx, user.ID)
}

func (u *Service) EnableTwoFA(ctx context.Context, userID, token string) ([]byte, error) {
	image, secret, err := u.otp.GenerateKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = u.db.EnableOTP(ctx, userID, secret)
	if err != nil {
		return nil, err
	}

	u.Logout(ctx, token)
	return image, nil
}

func (u *Service) DisableTwoFA(ctx context.Context, userID, token string) error {
	err := u.db.DisableOTP(ctx, userID)
	if err != nil {
		return err
	}

	u.Logout(ctx, token)
	return nil
}

func (u *Service) Logout(ctx context.Context, token string) {
	u.jwt.DeleteSalt(ctx, token)
}

func NewUserService(db DB, jwt *jwt.JWT, otp *otp.OTP, bcrypt *bcrypt.Bcrypt) *Service {
	return &Service{
		db:     db,
		jwt:    jwt,
		otp:    otp,
		bcrypt: bcrypt,
	}
}
