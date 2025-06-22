package auth

import (
	"lexia/ent"
	"lexia/internal/modules/user"
	"lexia/internal/shared"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getTokenPayloadDtoFromUserEntity(userEntity *ent.User) (*tokenPayloadDTO, error) {
	accessToken, err := shared.GenerateAccessToken(
		&shared.TokenClaims{
			UserID: userEntity.ID,
			Email:  userEntity.Email,
		},
	)

	if err != nil {
		return nil, err
	}

	return &tokenPayloadDTO{
		AccessToken: accessToken,
		User:        user.UserEntityToDto(userEntity),
	}, nil
}
