package auth

import (
	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/user"
	"github.com/kume1a/sonifybackend/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

func (dto *googleSignInDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *emailSignInDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getTokenPayloadDtoFromUserEntity(userEntity *database.User) (*tokenPayloadDTO, error) {
	accessToken, err := shared.GenerateAccessToken(&shared.TokenClaims{
		UserID: userEntity.ID,
		Email:  userEntity.Email.String})

	if err != nil {
		return nil, err
	}

	return &tokenPayloadDTO{
		AccessToken: accessToken,
		User:        user.UserEntityToDto(userEntity),
	}, nil
}
