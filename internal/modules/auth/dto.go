package auth

import "lexia/internal/modules/user"

type tokenPayloadDTO struct {
	AccessToken string       `json:"accessToken"`
	User        user.UserDto `json:"user"`
}

type emailSignInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailSignUpDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
