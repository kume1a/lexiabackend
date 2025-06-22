package auth

import (
	"time"

	"github.com/google/uuid"
)

type UserDto struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Username  string    `json:"name"`
	Email     string    `json:"email"`
}

type tokenPayloadDTO struct {
	AccessToken string  `json:"accessToken"`
	User        UserDto `json:"user"`
}

type emailSignInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
