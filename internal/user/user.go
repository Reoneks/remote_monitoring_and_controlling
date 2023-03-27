package user

import (
	"context"
	"project/pkg/bcrypt"
	"project/pkg/jwt"
	"project/pkg/otp"
	"project/structs"

	"github.com/google/uuid"
)

type UserService struct {
	db     DB
	jwt    *jwt.JWT
	otp    *otp.OTP
	bcrypt *bcrypt.Bcrypt
}

func (u *UserService) Register(ctx context.Context, req *structs.Register) ([]byte, error) {
	var (
		id     = uuid.NewString()
		image  []byte
		secret string
		err    error
	)

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
		Phone:      req.Phone,
		Password:   encryptedPassword,
		OTPEnabled: req.OTPEnabled,
		OTPSecret:  secret,
	})
}

func (u *UserService) Login(ctx context.Context, req *structs.Login) (string, bool, error) {
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

func (u *UserService) OTPCheck(ctx context.Context, req *structs.TwoFA) (string, error) {
	user, err := u.db.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		return "", err
	}

	ok := u.otp.ValidateKey(ctx, req.Password, user.OTPSecret)
	if !ok {
		return "", ErrInvalidOtpCode
	}

	return u.jwt.GenerateToken(ctx, user.ID)
}

func (u *UserService) EnableTwoFA(ctx context.Context, userID string) ([]byte, error) {
	image, secret, err := u.otp.GenerateKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	return image, u.db.EnableOTP(ctx, userID, secret)
}

func (u *UserService) Logout(ctx context.Context, token string) {
	u.jwt.ResetSalt(ctx, token)
	return
}

func NewUserService(db DB, jwt *jwt.JWT, otp *otp.OTP, bcrypt *bcrypt.Bcrypt) *UserService {
	return &UserService{
		db:     db,
		jwt:    jwt,
		otp:    otp,
		bcrypt: bcrypt,
	}
}
