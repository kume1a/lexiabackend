package word

import (
	"time"

	"github.com/google/uuid"
)

type CreateWordDTO struct {
	Text       string    `json:"text" validate:"required,min=1,max=500"`
	Definition string    `json:"definition" validate:"max=2000"`
	FolderID   uuid.UUID `json:"folderId" validate:"required"`
}

type UpdateWordDTO struct {
	Text       *string `json:"text" validate:"omitempty,min=1,max=500"`
	Definition *string `json:"definition" validate:"omitempty,max=2000"`
}

type WordDTO struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Text       string    `json:"text"`
	Definition string    `json:"definition"`
	FolderID   uuid.UUID `json:"folderId"`
}

type WordWithFolderDTO struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Text       string    `json:"text"`
	Definition string    `json:"definition"`
	Folder     struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	} `json:"folder"`
}

type WordDuplicateCheckDTO struct {
	IsDuplicate bool                   `json:"isDuplicate"`
	Word        *WordWithFolderPathDTO `json:"word,omitempty"`
}

type WordWithFolderPathDTO struct {
	ID         uuid.UUID           `json:"id"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
	Text       string              `json:"text"`
	Definition string              `json:"definition"`
	FolderPath []FolderPathItemDTO `json:"folderPath"`
}

type FolderPathItemDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
