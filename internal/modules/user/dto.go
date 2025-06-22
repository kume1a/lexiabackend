package user

import (
	"time"

	"github.com/google/uuid"
)

type updateUserDTO struct {
	Name string `json:"name" valid:"optional"`
}

type UserDto struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}
