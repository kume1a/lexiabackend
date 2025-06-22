package user

import (
	"lexia/ent"
)

func UserEntityToDto(userEntity *ent.User) UserDto {
	return UserDto{
		ID:        userEntity.ID,
		CreatedAt: userEntity.CreateTime,
		Name:      userEntity.Username,
		Email:     userEntity.Email,
	}
}
