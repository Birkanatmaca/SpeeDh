package model

import "time"

type User struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	FirstName           string     `gorm:"not null" json:"first_name"`
	LastName            string     `gorm:"not null" json:"last_name"`
	Email               string     `gorm:"unique;not null" json:"email"`
	Password            string     `gorm:"not null" json:"-"`
	IsActive            bool       `gorm:"default:false;not null" json:"is_active"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	VerificationCode    *string    `json:"-"`
	VerificationExpires *time.Time `json:"-"`
	VerificationSends   int        `gorm:"default:1;not null" json:"-"`
	LockedUntil         *time.Time `json:"-"`
}
