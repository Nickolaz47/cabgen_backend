package models

type AdminRegisterInput struct {
	RegisterInput
	UserRole UserRole `json:"user_role" binding:"required"`
}

type AdminUpdateInput struct {
	UpdateUserInput
	Email       *string  `json:"email,omitempty" binding:"omitempty,email"`
	Password    *string  `json:"password,omitempty" binding:"omitempty,min=8"`
	UserRole    *UserRole `json:"user_role,omitempty" binding:"omitempty"`
}
