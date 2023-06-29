package postgres

import (
	"context"

	"gorm.io/gorm/clause"
)

func (p *Postgres) GetUsers(ctx context.Context) ([]User, error) {
	var result []User
	return result, p.db.WithContext(ctx).Model(&User{}).Preload("ContactInfo").
		Select([]string{
			"users.id",
			"users.department",
			"users.position",
			"users.full_name",
		}).Find(&result).Error
}

func (p *Postgres) GetUserByPhone(ctx context.Context, phone string) (User, error) {
	var result User
	return result, p.db.WithContext(ctx).Model(&User{}).
		Select([]string{
			"users.id",
			"users.department",
			"users.position",
			"users.full_name",
			"users.password",
			"users.otp_secret",
		}).
		Joins("INNER JOIN contact_info ON users.id = contact_info.user_id").
		Where("contact_info.phone = ?", phone).
		First(&result).Error
}

func (p *Postgres) CreateUser(ctx context.Context, user *User) error {
	return p.db.Model(user).WithContext(ctx).Create(user).Error
}

func (p *Postgres) AddContactInfo(ctx context.Context, contactInfo *ContactInfo) error {
	return p.db.Model(contactInfo).WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "phone"}},
		DoUpdates: clause.AssignmentColumns([]string{"type"}),
	}).Create(contactInfo).Error
}

func (p *Postgres) EnableOTP(ctx context.Context, userID, secret string) error {
	return p.db.Model(&User{}).WithContext(ctx).Where("id = ?", userID).Update("otp_secret", secret).Error
}

func (p *Postgres) DisableOTP(ctx context.Context, userID string) error {
	return p.db.Model(&User{}).WithContext(ctx).Where("id = ?", userID).Update("otp_secret", "").Error
}

func (p *Postgres) DeleteUser(ctx context.Context, userID string) error {
	return p.db.Model(&User{}).WithContext(ctx).Where("id = ?", userID).Delete(&User{}).Error
}
