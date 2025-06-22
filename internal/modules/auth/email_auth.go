package auth

import (
	"context"
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
)

func AuthWithEmail(apiCfg config.ApiConfig, ctx context.Context, email string, password string) (*tokenPayloadDTO, *shared.HttpError) {
	userExistsByEmail, err := user.UserExistsByEmail(ctx, apiCfg.DB, email)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	var authUser *database.User
	var httpErr *shared.HttpError
	if userExistsByEmail {
		authUser, httpErr = signInWithEmail(apiCfg, ctx, email, password)
	} else {
		authUser, httpErr = signUpWithEmail(apiCfg, ctx, email, password)
	}

	if httpErr != nil {
		return nil, httpErr
	}

	tokenPayload, err := getTokenPayloadDtoFromUserEntity(authUser)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return tokenPayload, nil
}

func signInWithEmail(apiCfg config.ApiConfig, ctx context.Context, email string, password string) (*database.User, *shared.HttpError) {
	authUser, err := user.GetUserByEmail(ctx, apiCfg.DB, email)

	if err != nil {
		return nil, shared.Unauthorized(shared.ErrInvalidEmailOrPassword)
	}

	if authUser.AuthProvider != database.AuthProviderEMAIL {
		return nil, shared.BadRequest(shared.ErrInvalidAuthMethod)
	}

	if !ComparePasswordHash(password, authUser.PasswordHash.String) {
		return nil, shared.Unauthorized(shared.ErrInvalidEmailOrPassword)
	}

	return authUser, nil
}

func signUpWithEmail(apiCfg config.ApiConfig, ctx context.Context, email string, password string) (*database.User, *shared.HttpError) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	newUser, err := user.CreateUser(ctx, apiCfg.DB, database.CreateUserParams{
		Name:         sql.NullString{},
		Email:        sql.NullString{String: email, Valid: true},
		PasswordHash: sql.NullString{String: passwordHash, Valid: true},
		AuthProvider: database.AuthProviderEMAIL,
	})
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return newUser, nil
}
