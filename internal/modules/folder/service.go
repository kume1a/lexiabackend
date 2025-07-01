package folder

import (
	"context"
	"fmt"
	"lexia/ent"
	"lexia/ent/folder"
	"lexia/ent/schema"
	"lexia/ent/user"
	"lexia/internal/shared"

	"github.com/google/uuid"
)

type CreateFolderArgs struct {
	UserID       uuid.UUID
	Name         string
	Type         schema.FolderType
	LanguageFrom *schema.Language
	LanguageTo   *schema.Language
	ParentID     *uuid.UUID
}

type UpdateFolderArgs struct {
	FolderID uuid.UUID
	UserID   uuid.UUID
	Name     *string
	ParentID *uuid.UUID
}

func CreateFolder(ctx context.Context, db *ent.Client, args CreateFolderArgs) (*ent.Folder, error) {
	if args.Type == schema.FolderTypeWordCollection && args.LanguageFrom == nil {
		return nil, fmt.Errorf("languageFrom is required for word_collection folders")
	}
	if args.Type == schema.FolderTypeFolderCollection && (args.LanguageFrom != nil || args.LanguageTo != nil) {
		return nil, fmt.Errorf("languageFrom and languageTo should not be provided for folder collection folders")
	}

	if args.ParentID != nil {
		parentFolder, err := db.Folder.Query().
			Where(folder.ID(*args.ParentID)).
			WithUser().
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("parent folder not found: %w", err)
		}

		if parentFolder.Edges.User.ID != args.UserID {
			return nil, fmt.Errorf("parent folder does not belong to user")
		}

		if err := ValidateCanAddSubfolder(ctx, db, *args.ParentID); err != nil {
			return nil, err
		}
	}

	mutation := db.Folder.Create().
		SetName(args.Name).
		SetWordCount(0).
		SetType(args.Type).
		SetUserID(args.UserID)

	if args.Type == schema.FolderTypeWordCollection {
		if args.LanguageFrom != nil {
			mutation = mutation.SetLanguageFrom(*args.LanguageFrom)
		}
		if args.LanguageTo != nil {
			mutation = mutation.SetLanguageTo(*args.LanguageTo)
		}
	}

	if args.ParentID != nil {
		mutation = mutation.AddParentIDs(*args.ParentID)
	}

	return mutation.Save(ctx)
}

func GetFolderByID(ctx context.Context, db *ent.Client, folderID uuid.UUID, userID uuid.UUID) (*ent.Folder, error) {
	return db.Folder.Query().
		Where(folder.ID(folderID)).
		WithUser().
		WithWords().
		WithParent().
		WithSubfolders(func(q *ent.FolderQuery) {
			q.WithWords()
		}).
		Only(ctx)
}

func GetUserFolders(ctx context.Context, db *ent.Client, userID uuid.UUID) ([]*ent.Folder, error) {
	return db.Folder.Query().
		Where(folder.HasUserWith(user.ID(userID))).
		WithWords().
		WithParent().
		WithSubfolders(func(q *ent.FolderQuery) {
			q.WithWords()
		}).
		All(ctx)
}

func GetRootFolders(ctx context.Context, db *ent.Client, userID uuid.UUID) ([]*ent.Folder, error) {
	return db.Folder.Query().
		Where(
			folder.HasUserWith(user.ID(userID)),
			folder.Not(folder.HasParent()),
		).
		WithWords().
		WithSubfolders(func(q *ent.FolderQuery) {
			q.WithWords().
				WithSubfolders(func(q2 *ent.FolderQuery) {
					q2.WithWords()
				})
		}).
		All(ctx)
}

func UpdateFolder(ctx context.Context, db *ent.Client, args UpdateFolderArgs) (*ent.Folder, error) {
	existingFolder, err := db.Folder.Query().
		Where(folder.ID(args.FolderID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}

	if existingFolder.Edges.User.ID != args.UserID {
		return nil, fmt.Errorf("folder does not belong to user")
	}

	mutation := db.Folder.UpdateOneID(args.FolderID)

	if args.Name != nil {
		mutation = mutation.SetName(*args.Name)
	}

	if args.ParentID != nil {
		if err := validateParentChange(ctx, db, args.FolderID, *args.ParentID, args.UserID); err != nil {
			return nil, err
		}

		mutation = mutation.ClearParent().AddParentIDs(*args.ParentID)
	}

	return mutation.Save(ctx)
}

func DeleteFolder(ctx context.Context, db *ent.Client, folderID uuid.UUID, userID uuid.UUID) error {
	existingFolder, err := db.Folder.Query().
		Where(folder.ID(folderID)).
		WithUser().
		WithSubfolders().
		WithWords().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("folder not found: %w", err)
	}

	if existingFolder.Edges.User.ID != userID {
		return fmt.Errorf("folder does not belong to user")
	}

	if len(existingFolder.Edges.Subfolders) > 0 {
		return shared.BadRequest("Cannot delete folder that contains subfolders")
	}

	if len(existingFolder.Edges.Words) > 0 {
		return shared.BadRequest("Cannot delete folder that contains words")
	}

	return db.Folder.DeleteOneID(folderID).Exec(ctx)
}

func MoveFolder(ctx context.Context, db *ent.Client, folderID uuid.UUID, newParentID *uuid.UUID, userID uuid.UUID) (*ent.Folder, error) {
	existingFolder, err := db.Folder.Query().
		Where(folder.ID(folderID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}

	if existingFolder.Edges.User.ID != userID {
		return nil, fmt.Errorf("folder does not belong to user")
	}

	mutation := db.Folder.UpdateOneID(folderID).ClearParent()

	if newParentID != nil {
		if err := validateParentChange(ctx, db, folderID, *newParentID, userID); err != nil {
			return nil, err
		}
		mutation = mutation.AddParentIDs(*newParentID)
	}

	return mutation.Save(ctx)
}

func validateParentChange(ctx context.Context, db *ent.Client, folderID uuid.UUID, newParentID uuid.UUID, userID uuid.UUID) error {
	parentFolder, err := db.Folder.Query().
		Where(folder.ID(newParentID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("parent folder not found: %w", err)
	}

	if parentFolder.Edges.User.ID != userID {
		return fmt.Errorf("parent folder does not belong to user")
	}

	return checkCircularReference(ctx, db, folderID, newParentID)
}

func checkCircularReference(ctx context.Context, db *ent.Client, folderID uuid.UUID, potentialChildID uuid.UUID) error {
	if folderID == potentialChildID {
		return fmt.Errorf("cannot set folder as its own parent")
	}

	descendants, err := getDescendants(ctx, db, folderID)
	if err != nil {
		return err
	}

	for _, descendant := range descendants {
		if descendant == potentialChildID {
			return fmt.Errorf("circular reference detected: cannot move folder to its own descendant")
		}
	}

	return nil
}

func getDescendants(ctx context.Context, db *ent.Client, folderID uuid.UUID) ([]uuid.UUID, error) {
	var descendants []uuid.UUID

	subfolders, err := db.Folder.Query().
		Where(folder.HasParentWith(folder.ID(folderID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	for _, subfolder := range subfolders {
		descendants = append(descendants, subfolder.ID)
		subDescendants, err := getDescendants(ctx, db, subfolder.ID)
		if err != nil {
			return nil, err
		}
		descendants = append(descendants, subDescendants...)
	}

	return descendants, nil
}

func ValidateCanAddWords(ctx context.Context, db *ent.Client, folderID uuid.UUID) error {
	folder, err := db.Folder.Query().
		Where(folder.ID(folderID)).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("folder not found: %w", err)
	}

	if folder.Type != schema.FolderTypeWordCollection {
		return fmt.Errorf("words can only be added to word_collection folders")
	}

	return nil
}

func ValidateCanAddSubfolder(ctx context.Context, db *ent.Client, parentFolderID uuid.UUID) error {
	folder, err := db.Folder.Query().
		Where(folder.ID(parentFolderID)).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("parent folder not found: %w", err)
	}

	if folder.Type != schema.FolderTypeFolderCollection {
		return fmt.Errorf("subfolders can only be added to folder collection folders")
	}

	return nil
}

func GetFolderType(ctx context.Context, db *ent.Client, folderID uuid.UUID) (string, error) {
	folder, err := db.Folder.Query().
		Where(folder.ID(folderID)).
		Only(ctx)
	if err != nil {
		return "", fmt.Errorf("folder not found: %w", err)
	}

	return string(folder.Type), nil
}
