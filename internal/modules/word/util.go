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

func buildFolderPath(folder *ent.Folder) []FolderPathItemDTO {
	var path []FolderPathItemDTO
	current := folder

	for current != nil {
		path = append([]FolderPathItemDTO{{
			ID:   current.ID,
			Name: current.Name,
		}}, path...)

		if current.Edges.Parent != nil && len(current.Edges.Parent) > 0 {
			current = current.Edges.Parent[0]
		} else {
			current = nil
		}
	}

	return path
}

func WordEntityWithFolderPathToDTO(wordEntity *ent.Word) WordWithFolderPathDTO {
	dto := WordWithFolderPathDTO{
		ID:         wordEntity.ID,
		CreatedAt:  wordEntity.CreateTime,
		UpdatedAt:  wordEntity.UpdateTime,
		Text:       wordEntity.Text,
		Definition: wordEntity.Definition,
	}

	if wordEntity.Edges.Folder != nil {
		dto.FolderPath = buildFolderPath(wordEntity.Edges.Folder)
	}

	return dto
}
