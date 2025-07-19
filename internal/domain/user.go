package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required,min=3,max=25"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
