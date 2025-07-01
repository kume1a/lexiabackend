package folder

import (
	"lexia/ent/schema"
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleCreateFolder(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		var body CreateFolderDTO
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
			return
		}

		if body.Type == schema.FolderTypeWordCollection && body.LanguageFrom == nil {
			shared.ResBadRequest(c, "languageFrom is required for word_collection folders")
			return
		}
		if body.Type == schema.FolderTypeFolderCollection && (body.LanguageFrom != nil || body.LanguageTo != nil) {
			shared.ResBadRequest(c, "languageFrom and languageTo should not be provided for folder collection folders")
			return
		}

		folder, err := CreateFolder(
			c.Request.Context(), apiCfg.DB,
			CreateFolderArgs{
				UserID:       authPayload.UserID,
				Name:         body.Name,
				Type:         body.Type,
				LanguageFrom: body.LanguageFrom,
				LanguageTo:   body.LanguageTo,
				ParentID:     body.ParentID,
			},
		)

		if err != nil {
			shared.ResBadRequest(c, err.Error())
			return
		}

		shared.ResOK(c, FolderEntityToDto(folder))
	}
}

func handleGetFolder(apiCfg *shared.ApiConfig) gin.HandlerFunc {
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

		folder, err := GetFolderByID(c.Request.Context(), apiCfg.DB, folderID, authPayload.UserID)
		if err != nil {
			shared.ResNotFound(c, "Folder not found")
			return
		}

		shared.ResOK(c, FolderEntityToDto(folder))
	}
}

func handleGetUserFolders(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		folders, err := GetUserFolders(c.Request.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		folderDTOs := make([]FolderDTO, len(folders))
		for i, folder := range folders {
			folderDTOs[i] = FolderEntityToDto(folder)
		}

		shared.ResOK(c, folderDTOs)
	}
}

func handleGetRootFolders(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		folders, err := GetRootFolders(c.Request.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		folderDTOs := make([]FolderDTO, len(folders))
		for i, folder := range folders {
			folderDTOs[i] = FolderEntityToDto(folder)
		}

		shared.ResOK(c, folderDTOs)
	}
}

func handleUpdateFolder(apiCfg *shared.ApiConfig) gin.HandlerFunc {
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

		var body UpdateFolderDTO
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
			return
		}

		folder, err := UpdateFolder(
			c.Request.Context(), apiCfg.DB,
			UpdateFolderArgs{
				FolderID: folderID,
				UserID:   authPayload.UserID,
				Name:     body.Name,
				ParentID: body.ParentID,
			},
		)

		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		shared.ResOK(c, FolderEntityToDto(folder))
	}
}

func handleDeleteFolder(apiCfg *shared.ApiConfig) gin.HandlerFunc {
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

		err = DeleteFolder(c.Request.Context(), apiCfg.DB, folderID, authPayload.UserID)
		if err != nil {
			if httpErr, ok := err.(*shared.HttpError); ok {
				shared.ResHttpError(c, httpErr)
				return
			}
			shared.ResInternalServerErrorDef(c)
			return
		}

		shared.ResNoContent(c)
	}
}

func handleMoveFolder(apiCfg *shared.ApiConfig) gin.HandlerFunc {
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

		type MoveFolderDTO struct {
			ParentID *uuid.UUID `json:"parentId"`
		}

		var body MoveFolderDTO
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
			return
		}

		folder, err := MoveFolder(c.Request.Context(), apiCfg.DB, folderID, body.ParentID, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		shared.ResOK(c, FolderEntityToDto(folder))
	}
}
