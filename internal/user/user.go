package user

import (
	"context"
	"time"

	"remote_monitoring_and_controlling/pkg/cache"

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
	cache  *cache.Cache[structs.User]
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

	phoneNumber = "+" + phoneNumber
	return image, u.db.CreateUser(ctx, &structs.User{
		ID:         id,
		Phone:      phoneNumber,
		Password:   encryptedPassword,
		OTPEnabled: req.OTPEnabled,
		OTPSecret:  secret,
	})
}

func (u *Service) Login(ctx context.Context, req *structs.Login) (string, bool, error) {
	phoneNumber := phonenumber.Parse(req.Phone, "")
	if phoneNumber == "" {
		return "", false, ErrInvalidPhone
	}

	phoneNumber = "+" + phoneNumber
	user, err := u.db.GetUserByPhone(ctx, phoneNumber)
	if err != nil {
		return "", false, err
	}

	if !u.bcrypt.Validate(ctx, user.Password, req.Password) {
		return "", false, ErrInvalidPassword
	}

	if user.OTPEnabled {
		id := ulid.Make().String()
		u.cache.Set(ctx, id, user, 3*time.Minute)
		return id, true, nil
	}

	token, err := u.jwt.GenerateToken(ctx, user.ID)
	return token, false, err
}

func (u *Service) AddAlternativeNumber(ctx context.Context, userID, phone string) error {
	phoneNumber := phonenumber.Parse(phone, "")
	if phoneNumber == "" {
		return ErrInvalidPhone
	}

	phoneNumber = "+" + phoneNumber
	user, err := u.db.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.ID = ulid.Make().String()
	user.Phone = phoneNumber
	return u.db.CreateUser(ctx, &user)
}

func (u *Service) OTPCheck(ctx context.Context, req *structs.TwoFA) (string, error) {
	user, found := u.cache.Get(ctx, req.ID)
	if !found {
		return "", ErrInvalid2FAID
	}

	u.cache.Delete(ctx, req.ID)
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
		cache:  cache.NewCache[structs.User](),
	}
}
