package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"not null;uniqueIndex;size:20"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserRole struct {
	UserID int `json:"user_id" gorm:"primaryKey"`
	RoleID int `json:"role_id" gorm:"primaryKey"`
}
