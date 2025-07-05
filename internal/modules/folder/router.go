package folder

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func Router(apiCfg *shared.ApiConfig, rg *gin.RouterGroup) {
	folderGroup := rg.Group("/folders")
	{
		folderGroup.GET("", handleGetUserFolders(apiCfg))
		folderGroup.GET("/root", handleGetRootFolders(apiCfg))
		folderGroup.GET("/:folderId", handleGetFolder(apiCfg))
		folderGroup.GET("/:folderId/subfolders", handleGetSubfoldersByFolderID(apiCfg))

		folderGroup.POST("", handleCreateFolder(apiCfg))

		folderGroup.PUT("/:folderId", handleUpdateFolder(apiCfg))
		folderGroup.PUT("/:folderId/move", handleMoveFolder(apiCfg))

		folderGroup.DELETE("/:folderId", handleDeleteFolder(apiCfg))
	}
}
