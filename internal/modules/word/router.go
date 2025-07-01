package word

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func Router(apiCfg *shared.ApiConfig, rg *gin.RouterGroup) {
	wordGroup := rg.Group("/words")
	{
		wordGroup.POST("", handleCreateWord(apiCfg))
		wordGroup.GET("/:wordId", handleGetWord(apiCfg))
		wordGroup.PUT("/:wordId", handleUpdateWord(apiCfg))
		wordGroup.DELETE("/:wordId", handleDeleteWord(apiCfg))
	}

	folderGroup := rg.Group("/folders")
	{
		folderGroup.GET("/:folderId/words", handleGetWordsByFolder(apiCfg))
	}
}
