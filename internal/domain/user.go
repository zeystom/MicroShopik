package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	ID        int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string         `json:"username" gorm:"not null;uniqueIndex:uni_users_username;size:100"`
	Email     string         `json:"email" gorm:"not null;uniqueIndex:uni_users_email;size:255"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	Roles     []Role         `json:"roles" gorm:"many2many:user_roles;"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	LastLogin *time.Time     `json:"last_login"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type JWTClaims struct {
	UserID int      `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}
