package user

import (
	"context"
	"log"

	"entgo.io/ent"
	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, client *ent.Client, username, email, password string) (*ent.User, error) {
	newUser, err := client.User.Create().
		SetID(uuid.New()).
		SetUsername(username).
		SetEmail(email).
		SetPassword(password).
		Save(ctx)

	if err != nil {
		log.Println("Error creating user: ", err)
		return nil, err
	}

	return newUser, nil
}

func GetUserByID(ctx context.Context, client *ent.Client, userID uuid.UUID) (*ent.User, error) {
	user, err := client.User.Get(ctx, userID)

	if ent.IsNotFound(err) {
		return nil, shared.NotFound(shared.ErrUserNotFound)
	}

	if err != nil {
		log.Println("Error getting user by ID: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return user, nil
}

func GetUserByEmail(ctx context.Context, client *ent.Client, email string) (*ent.User, error) {
	user, err := client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, shared.NotFound(shared.ErrUserNotFound)
	}

	if err != nil {
		log.Println("Error getting user by email: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return user, nil
}

func UpdateUser(ctx context.Context, client *ent.Client, userID uuid.UUID, username string) (*ent.User, error) {
	updatedUser, err := client.User.UpdateOneID(userID).
		SetUsername(username).
		Save(ctx)

	if err != nil {
		log.Println("Error updating user: ", err)
		return nil, err
	}

	return updatedUser, nil
}

func UserExistsByEmail(ctx context.Context, client *ent.Client, email string) (bool, error) {
	count, err := client.User.Query().
		Where(user.EmailEQ(email)).
		Count(ctx)

	if err != nil {
		log.Println("Error counting users by email: ", err)
		return false, err
	}

	return count > 0, nil
}
