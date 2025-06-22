package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func CreateUser(
	ctx context.Context,
	db *database.Queries,
	params database.CreateUserParams,
) (*database.User, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	user, err := db.CreateUser(ctx, params)

	if err != nil {
		log.Println("Error creating user: ", err)
	}

	return &user, err
}

func GetUserByID(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (*database.User, error) {
	user, err := db.GetUserByID(ctx, userId)

	if shared.IsDBErrorNotFound(err) {
		return nil, shared.NotFound(shared.ErrUserNotFound)
	}

	if err != nil {
		log.Println("Error getting user by ID: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &user, nil
}

func GetUserByEmail(ctx context.Context, db *database.Queries, email string) (*database.User, error) {
	user, err := db.GetUserByEmail(ctx, sql.NullString{
		String: email,
		Valid:  true,
	})

	if err != nil {
		log.Println("Error getting user by email: ", err)
	}

	return &user, err
}

func UpdateUser(
	ctx context.Context,
	db *database.Queries,
	params *database.UpdateUserParams,
) (*database.User, error) {
	user, err := db.UpdateUser(ctx, *params)

	if err != nil {
		log.Println("Error updating user: ", err)
	}

	return &user, err
}

func UserExistsByEmail(
	ctx context.Context,
	db *database.Queries,
	email string,
) (bool, error) {
	count, err := db.CountUsersByEmail(ctx, sql.NullString{String: email, Valid: true})

	if err != nil {
		log.Println("Error counting users by email: ", err)
	}

	return count > 0, err
}
