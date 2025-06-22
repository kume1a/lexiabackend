package user

import (
	"lexia/ent/schema"
	"lexia/internal/modules/user"
)

func UserEntityToDto(userEntity *schema.User) user.UserDto {
	return user.UserDto{
		ID:           userEntity.,
		CreatedAt:    userEntity.CreatedAt,
		Name:         userEntity.Name.String,
		Email:        userEntity.Email.String,
		AuthProvider: userEntity.AuthProvider,
	}
}
