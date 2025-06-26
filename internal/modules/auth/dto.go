package auth

import "lexia/internal/modules/user"

type tokenPayloadDTO struct {
	AccessToken string       `json:"accessToken"`
	User        user.UserDto `json:"user"`
}

type emailSignInDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type EmailSignUpDTO struct {
	Username string `json:"username" validate:"username" binding:"required,min=2,max=50,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" validate:"strong_password" binding:"required,min=8,max=128"`
}
