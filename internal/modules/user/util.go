package user

import (
	"lexia/ent"
)

func UserEntityToDto(userEntity *ent.User) UserDto {
	return UserDto{
		ID:        userEntity.ID,
		CreatedAt: userEntity.CreateTime,
		Username:  userEntity.Username,
		Email:     userEntity.Email,
	}
}
