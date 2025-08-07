package models

type AdminRegisterInput struct {
	RegisterInput
	UserRole UserRole `json:"user_role" binding:"required"`
}

