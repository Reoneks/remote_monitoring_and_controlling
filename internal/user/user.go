package user

import (
	"context"
	"time"

	"remote_monitoring_and_controlling/pkg/cache"
	"remote_monitoring_and_controlling/pkg/postgres"

	"remote_monitoring_and_controlling/pkg/bcrypt"
	"remote_monitoring_and_controlling/pkg/jwt"
	"remote_monitoring_and_controlling/pkg/otp"

	"github.com/dongri/phonenumber"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Service struct {
	db     DB
	jwt    *jwt.JWT
	otp    *otp.OTP
	bcrypt *bcrypt.Bcrypt
	cache  *cache.Cache[OTPData]
}

func (u *Service) Register(ctx context.Context, req *Register) error {
	id := ulid.Make().String()
	contactInfo := make([]postgres.ContactInfo, 0, len(req.ContactInfo))
	for _, info := range req.ContactInfo {
		phoneNumber := phonenumber.Parse(info.Phone, "")
		if phoneNumber == "" {
			return ErrInvalidPhone
		}

		contactInfo = append(contactInfo, postgres.ContactInfo{
			UserID: id,
			Type:   info.Type,
			Phone:  "+" + phoneNumber,
		})
	}

	encryptedPassword, err := u.bcrypt.Encode(ctx, req.Password)
	if err != nil {
		return err
	}

	err = u.db.CreateUser(ctx, &postgres.User{
		ID:          id,
		Department:  req.Department,
		Position:    req.Position,
		FullName:    req.FullName,
		ForeignID:   req.ForeignID,
		Password:    encryptedPassword,
		ContactInfo: contactInfo,
	})
	if err != nil {
		return err
	}

	for i := range contactInfo {
		err = u.db.AddContactInfo(ctx, &contactInfo[i])
		if err != nil {
			if err := u.db.DeleteUser(ctx, id); err != nil {
				log.Error().Str("function", "Register").Err(err).Msg("Failed to delete user")
			}

			return err
		}
	}

	return nil
}

func (u *Service) Login(ctx context.Context, req *Login) (string, bool, error) {
	phoneNumber := phonenumber.Parse(req.Phone, "")
	if phoneNumber == "" {
		return "", false, ErrInvalidPhone
	}

	user, err := u.db.GetUserByPhone(ctx, "+"+phoneNumber)
	if err != nil {
		return "", false, err
	}

	if !u.bcrypt.Validate(ctx, user.Password, req.Password) {
		return "", false, ErrInvalidPassword
	}

	if user.OTPSecret != "" {
		id := ulid.Make().String()
		u.cache.Set(ctx, id, OTPData{UserID: user.ID, OTPSecret: user.OTPSecret}, 3*time.Minute)
		return id, true, nil
	}

	token, err := u.jwt.GenerateToken(ctx, user.ID)
	return token, false, err
}

func (u *Service) AddAlternativeNumber(ctx context.Context, req *AddAlternativeNumber) error {
	id, err := u.db.GetUserIDByForeignID(ctx, req.ForeignID)
	if err != nil {
		return err
	}

	for _, info := range req.ContactInfo {
		phoneNumber := phonenumber.Parse(info.Phone, "")
		if phoneNumber == "" {
			return ErrInvalidPhone
		}

		err := u.db.AddContactInfo(ctx, &postgres.ContactInfo{
			UserID: id,
			Type:   info.Type,
			Phone:  "+" + phoneNumber,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Service) OTPCheck(ctx context.Context, req *TwoFA) (string, error) {
	data, found := u.cache.Get(ctx, req.ID)
	if !found {
		return "", ErrInvalid2FAID
	}

	u.cache.Delete(ctx, req.ID)
	ok := u.otp.ValidateKey(ctx, req.OTPPassword, data.OTPSecret)
	if !ok {
		return "", ErrInvalidOtpCode
	}

	return u.jwt.GenerateToken(ctx, data.UserID)
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
		cache:  cache.NewCache[OTPData](),
	}
}
