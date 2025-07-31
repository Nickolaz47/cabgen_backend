package models

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserToken struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	UserRole UserRole  `json:"user_role"`
	jwt.RegisteredClaims
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
