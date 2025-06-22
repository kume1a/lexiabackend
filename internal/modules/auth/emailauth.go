package auth

import (
	"context"
	"lexia/internal/modules/user"
	"lexia/internal/shared"
)

type SignInWithEmailArgs struct {
	Email    string
	Password string
}

func SignInWithEmail(
	apiCfg *shared.ApiConfig,
	ctx context.Context,
	args SignInWithEmailArgs,
) (*tokenPayloadDTO, *shared.HttpError) {
	authUser, err := user.GetUserByEmail(ctx, apiCfg.DB, args.Email)

	if err != nil {
		if shared.IsDatabaseErorNotFound(err) {
			return nil, shared.Unauthorized(shared.ErrInvalidEmailOrPassword)
		}

		return nil, shared.InternalServerErrorDef()
	}

	if !ComparePasswordHash(args.Password, authUser.Password) {
		return nil, shared.Unauthorized(shared.ErrInvalidEmailOrPassword)
	}

	tokenPayload, err := getTokenPayloadDtoFromUserEntity(authUser)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return tokenPayload, nil
}

type SignUpWithEmailArgs struct {
	Username string
	Email    string
	Password string
}

func SignUpWithEmail(
	apiCfg *shared.ApiConfig,
	ctx context.Context,
	args SignUpWithEmailArgs,
) (*tokenPayloadDTO, *shared.HttpError) {
	userExistsByEmail, err := user.UserExistsByEmail(ctx, apiCfg.DB, args.Email)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	if userExistsByEmail {
		return nil, shared.BadRequest(shared.ErrEmailAlreadyExists)
	}

	passwordHash, err := HashPassword(args.Password)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	newUser, err := user.CreateUser(ctx, apiCfg.DB, user.CreateUserArgs{
		Username: args.Username,
		Email:    args.Email,
		Password: passwordHash,
	})
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	tokenPayload, err := getTokenPayloadDtoFromUserEntity(newUser)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return tokenPayload, nil
}
