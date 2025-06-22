package user

import (
	"context"
	"lexia/ent"
	"lexia/ent/user"
	"log"

	"github.com/google/uuid"
)

type CreateUserArgs struct {
	Username string
	Email    string
	Password string
}

func CreateUser(
	ctx context.Context,
	db *ent.Client,
	args CreateUserArgs,
) (*ent.User, error) {
	newUser, err := db.User.Create().
		SetID(uuid.New()).
		SetUsername(args.Username).
		SetEmail(args.Email).
		SetPassword(args.Password).
		Save(ctx)

	if err != nil {
		log.Println("Error creating user: ", err)
		return nil, err
	}

	return newUser, nil
}

func GetUserByID(
	ctx context.Context,
	db *ent.Client,
	ID uuid.UUID,
) (*ent.User, error) {
	user, err := db.User.Get(ctx, ID)

	if ent.IsNotFound(err) {
		log.Println("User not found with ID: ", ID)
		return nil, err
	}

	if err != nil {
		log.Println("Error getting user by ID: ", err)
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(
	ctx context.Context,
	db *ent.Client,
	email string,
) (*ent.User, error) {
	user, err := db.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if err != nil {
		log.Println("Error getting user by email: ", err)
		return nil, err
	}

	return user, nil
}

type UpdateUserByIDArgs struct {
	UserID   uuid.UUID
	Username string
}

func UpdateUserByID(
	ctx context.Context,
	db *ent.Client,
	args UpdateUserByIDArgs,
) (*ent.User, error) {
	updatedUser, err := db.User.UpdateOneID(args.UserID).
		SetUsername(args.Username).
		Save(ctx)

	if err != nil {
		log.Println("Error updating user: ", err)
		return nil, err
	}

	return updatedUser, nil
}

func UserExistsByEmail(
	ctx context.Context,
	db *ent.Client,
	email string,
) (bool, error) {
	count, err := db.User.Query().
		Where(user.EmailEQ(email)).
		Count(ctx)

	if err != nil {
		log.Println("Error counting users by email: ", err)
		return false, err
	}

	return count > 0, nil
}
