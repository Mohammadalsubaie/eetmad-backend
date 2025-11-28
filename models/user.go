package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Email      string         `gorm:"size:255;unique;not null" json:"email"`
	Phone      string         `gorm:"size:20;unique" json:"phone,omitempty"`
	Password   string         `gorm:"size:255;not null" json:"-"`
	UserType   string         `gorm:"size:20;default:'client'" json:"user_type"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	IsVerified bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
