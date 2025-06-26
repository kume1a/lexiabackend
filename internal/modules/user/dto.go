package user

import (
	"time"

	"github.com/google/uuid"
)

type updateUserDTO struct {
	Username string `json:"username" validate:"username" binding:"required,min=2,max=50,alphanum"`
}

type UserDto struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}
