package auth

import (
	"mainbackend/internal/model"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	ActivateUser(user *model.User) error
	UpdatePassword(userID uint, newPasswordHash string) error
	SetVerificationCode(user *model.User) error
	IncrementVerificationSends(user *model.User) error
	LockUserAccount(user *model.User, until time.Time) error
	DeleteUserByID(userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func (r *userRepository) ActivateUser(user *model.User) error {
	return r.db.Model(user).Updates(map[string]interface{}{
		"is_active":            true,
		"verification_code":    nil,
		"verification_expires": nil,
	}).Error
}

func (r *userRepository) UpdatePassword(userID uint, newPasswordHash string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password":             newPasswordHash,
		"verification_code":    nil,
		"verification_expires": nil,
	}).Error
}

func (r *userRepository) SetVerificationCode(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) IncrementVerificationSends(user *model.User) error {
	return r.db.Model(user).Update("verification_sends", gorm.Expr("verification_sends + 1")).Error
}

func (r *userRepository) LockUserAccount(user *model.User, until time.Time) error {
	return r.db.Model(user).Update("locked_until", until).Error
}

func (r *userRepository) DeleteUserByID(userID uint) error {
	return r.db.Where("id = ?", userID).Delete(&model.User{}).Error
}
