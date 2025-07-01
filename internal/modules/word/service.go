package word

import (
	"context"
	"fmt"
	"lexia/ent"
	"lexia/ent/folder"
	"lexia/ent/word"
	"log"

	"github.com/google/uuid"
)

type CreateWordArgs struct {
	Text       string
	Definition string
	FolderID   uuid.UUID
	UserID     uuid.UUID
}

type UpdateWordArgs struct {
	WordID     uuid.UUID
	UserID     uuid.UUID
	Text       *string
	Definition *string
}

func CreateWord(
	ctx context.Context,
	db *ent.Client,
	args CreateWordArgs,
) (*ent.Word, error) {
	folderEntity, err := db.Folder.Query().
		Where(folder.ID(args.FolderID)).
		WithUser().
		Only(ctx)
	if err != nil {
		log.Println("Error finding folder: ", err)
		return nil, err
	}

	if folderEntity.Edges.User.ID != args.UserID {
		return nil, fmt.Errorf("folder does not belong to user")
	}

	newWord, err := db.Word.Create().
		SetID(uuid.New()).
		SetText(args.Text).
		SetDefinition(args.Definition).
		SetFolderID(args.FolderID).
		Save(ctx)

	if err != nil {
		log.Println("Error creating word: ", err)
		return nil, err
	}

	return newWord, nil
}

func GetWordByID(
	ctx context.Context,
	db *ent.Client,
	wordID uuid.UUID,
) (*ent.Word, error) {
	word, err := db.Word.Get(ctx, wordID)

	if ent.IsNotFound(err) {
		log.Println("Word not found with ID: ", wordID)
		return nil, err
	}

	if err != nil {
		log.Println("Error getting word by ID: ", err)
		return nil, err
	}

	return word, nil
}

func GetWordByIDWithFolder(
	ctx context.Context,
	db *ent.Client,
	wordID uuid.UUID,
) (*ent.Word, error) {
	word, err := db.Word.Query().
		Where(word.ID(wordID)).
		WithFolder().
		Only(ctx)

	if ent.IsNotFound(err) {
		log.Println("Word not found with ID: ", wordID)
		return nil, err
	}

	if err != nil {
		log.Println("Error getting word by ID with folder: ", err)
		return nil, err
	}

	return word, nil
}

func GetWordsByFolderID(
	ctx context.Context,
	db *ent.Client,
	folderID uuid.UUID,
	userID uuid.UUID,
) ([]*ent.Word, error) {
	folderEntity, err := db.Folder.Query().
		Where(folder.ID(folderID)).
		WithUser().
		Only(ctx)
	if err != nil {
		log.Println("Error finding folder: ", err)
		return nil, err
	}

	if folderEntity.Edges.User.ID != userID {
		return nil, fmt.Errorf("folder does not belong to user")
	}

	words, err := db.Word.Query().
		Where(word.HasFolderWith(folder.ID(folderID))).
		WithFolder().
		All(ctx)

	if err != nil {
		log.Println("Error getting words by folder ID: ", err)
		return nil, err
	}

	return words, nil
}

func UpdateWord(
	ctx context.Context,
	db *ent.Client,
	args UpdateWordArgs,
) (*ent.Word, error) {
	wordEntity, err := db.Word.Query().
		Where(word.ID(args.WordID)).
		WithFolder(func(q *ent.FolderQuery) {
			q.WithUser()
		}).
		Only(ctx)
	if err != nil {
		log.Println("Error finding word: ", err)
		return nil, err
	}

	if wordEntity.Edges.Folder.Edges.User.ID != args.UserID {
		return nil, fmt.Errorf("word does not belong to user")
	}

	updateQuery := db.Word.UpdateOneID(args.WordID)

	if args.Text != nil {
		updateQuery = updateQuery.SetText(*args.Text)
	}

	if args.Definition != nil {
		updateQuery = updateQuery.SetDefinition(*args.Definition)
	}

	updatedWord, err := updateQuery.Save(ctx)

	if err != nil {
		log.Println("Error updating word: ", err)
		return nil, err
	}

	return updatedWord, nil
}

func DeleteWord(
	ctx context.Context,
	db *ent.Client,
	wordID uuid.UUID,
	userID uuid.UUID,
) error {
	wordEntity, err := db.Word.Query().
		Where(word.ID(wordID)).
		WithFolder(func(q *ent.FolderQuery) {
			q.WithUser()
		}).
		Only(ctx)
	if err != nil {
		log.Println("Error finding word: ", err)
		return err
	}

	if wordEntity.Edges.Folder.Edges.User.ID != userID {
		return fmt.Errorf("word does not belong to user")
	}

	err = db.Word.DeleteOneID(wordID).Exec(ctx)

	if err != nil {
		log.Println("Error deleting word: ", err)
		return err
	}

	return nil
}
