package folder

import (
	"lexia/ent"
	"lexia/ent/schema"

	"github.com/google/uuid"
)

type CreateFolderDTO struct {
	Name         string            `json:"name" validate:"required,min=1,max=255"`
	Type         schema.FolderType `json:"type" validate:"required,oneof=FOLDER_COLLECTION WORD_COLLECTION"`
	LanguageFrom *schema.Language  `json:"languageFrom,omitempty"`
	LanguageTo   *schema.Language  `json:"languageTo,omitempty"`
	ParentID     *uuid.UUID        `json:"parentId,omitempty"`
}

type UpdateFolderDTO struct {
	Name     *string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	ParentID *uuid.UUID `json:"parentId,omitempty"`
}

type FolderDTO struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	Type         schema.FolderType `json:"type"`
	WordCount    int32             `json:"wordCount"`
	LanguageFrom *schema.Language  `json:"languageFrom,omitempty"`
	LanguageTo   *schema.Language  `json:"languageTo,omitempty"`
	ParentID     *uuid.UUID        `json:"parentId,omitempty"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
	Subfolders   []FolderDTO       `json:"subfolders,omitempty"`
	HasWords     bool              `json:"hasWords"`
}

func FolderEntityToDto(folder *ent.Folder) FolderDTO {
	dto := FolderDTO{
		ID:        folder.ID,
		Name:      folder.Name,
		Type:      folder.Type,
		WordCount: folder.WordCount,
		CreatedAt: folder.CreateTime.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: folder.UpdateTime.Format("2006-01-02T15:04:05Z"),
		HasWords:  len(folder.Edges.Words) > 0,
	}

	if folder.LanguageFrom != nil {
		dto.LanguageFrom = folder.LanguageFrom
	}

	if folder.LanguageTo != nil {
		dto.LanguageTo = folder.LanguageTo
	}

	if len(folder.Edges.Parent) > 0 {
		dto.ParentID = &folder.Edges.Parent[0].ID
	}

	if len(folder.Edges.Subfolders) > 0 {
		dto.Subfolders = make([]FolderDTO, len(folder.Edges.Subfolders))
		for i, subfolder := range folder.Edges.Subfolders {
			dto.Subfolders[i] = FolderEntityToDto(subfolder)
		}
	}

	return dto
}
