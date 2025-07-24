package word

import (
	"lexia/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleCreateWord(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		var body CreateWordDTO
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
			return
		}

		word, err := CreateWord(
			c.Request.Context(),
			apiCfg.DB,
			CreateWordArgs{
				Text:       body.Text,
				Definition: body.Definition,
				FolderID:   body.FolderID,
				UserID:     authPayload.UserID,
			},
		)

		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		wordWithFolder, err := GetWordByIDWithFolder(c.Request.Context(), apiCfg.DB, word.ID)
		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		shared.ResCreated(c, WordEntityWithFolderToDTO(wordWithFolder))
	}
}

func handleGetWord(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		wordIDStr := c.Param("wordId")
		wordID, err := uuid.Parse(wordIDStr)
		if err != nil {
			shared.ResBadRequest(c, "Invalid word ID")
			return
		}

		word, err := GetWordByIDWithFolder(c.Request.Context(), apiCfg.DB, wordID)
		if err != nil {
			shared.ResNotFound(c, "Word not found")
			return
		}

		shared.ResOK(c, WordEntityWithFolderToDTO(word))
	}
}

func handleGetWordsByFolder(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		folderIDStr := c.Param("folderId")
		folderID, err := uuid.Parse(folderIDStr)
		if err != nil {
			shared.ResBadRequest(c, "Invalid folder ID")
			return
		}

		words, err := GetWordsByFolderID(
			c.Request.Context(),
			apiCfg.DB,
			folderID,
			authPayload.UserID,
		)

		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		shared.ResOK(c, WordEntitiesToDTOs(words))
	}
}

func handleUpdateWord(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		wordIDStr := c.Param("wordId")
		wordID, err := uuid.Parse(wordIDStr)
		if err != nil {
			shared.ResBadRequest(c, "Invalid word ID")
			return
		}

		var body UpdateWordDTO
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
			return
		}

		word, err := UpdateWord(
			c.Request.Context(),
			apiCfg.DB,
			UpdateWordArgs{
				WordID:     wordID,
				UserID:     authPayload.UserID,
				Text:       body.Text,
				Definition: body.Definition,
			},
		)

		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		wordWithFolder, err := GetWordByIDWithFolder(c.Request.Context(), apiCfg.DB, word.ID)
		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		shared.ResOK(c, WordEntityWithFolderToDTO(wordWithFolder))
	}
}

func handleDeleteWord(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		wordIDStr := c.Param("wordId")
		wordID, err := uuid.Parse(wordIDStr)
		if err != nil {
			shared.ResBadRequest(c, "Invalid word ID")
			return
		}

		err = DeleteWord(
			c.Request.Context(),
			apiCfg.DB,
			wordID,
			authPayload.UserID,
		)

		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func handleCheckWordDuplicate(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		text := c.Query("text")
		if text == "" {
			shared.ResBadRequest(c, "Text parameter is required")
			return
		}

		duplicateWord, err := CheckWordDuplicate(
			c.Request.Context(),
			apiCfg.DB,
			text,
			authPayload.UserID,
		)

		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		response := WordDuplicateCheckDTO{
			IsDuplicate: duplicateWord != nil,
		}

		if duplicateWord != nil {
			wordDTO := WordEntityWithFolderPathToDTO(duplicateWord)
			response.Word = &wordDTO
		}

		shared.ResOK(c, response)
	}
}
