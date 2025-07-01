package word

import (
	"lexia/ent"

	"github.com/google/uuid"
)

func WordEntityToDTO(wordEntity *ent.Word) WordDTO {
	var folderID uuid.UUID
	if wordEntity.Edges.Folder != nil {
		folderID = wordEntity.Edges.Folder.ID
	}

	return WordDTO{
		ID:         wordEntity.ID,
		CreatedAt:  wordEntity.CreateTime,
		UpdatedAt:  wordEntity.UpdateTime,
		Text:       wordEntity.Text,
		Definition: wordEntity.Definition,
		FolderID:   folderID,
	}
}

func WordEntityWithFolderToDTO(wordEntity *ent.Word) WordWithFolderDTO {
	dto := WordWithFolderDTO{
		ID:         wordEntity.ID,
		CreatedAt:  wordEntity.CreateTime,
		UpdatedAt:  wordEntity.UpdateTime,
		Text:       wordEntity.Text,
		Definition: wordEntity.Definition,
	}

	if wordEntity.Edges.Folder != nil {
		dto.Folder.ID = wordEntity.Edges.Folder.ID
		dto.Folder.Name = wordEntity.Edges.Folder.Name
	}

	return dto
}

func WordEntitiesToDTOs(wordEntities []*ent.Word) []WordDTO {
	dtos := make([]WordDTO, len(wordEntities))
	for i, wordEntity := range wordEntities {
		dtos[i] = WordEntityToDTO(wordEntity)
	}
	return dtos
}
