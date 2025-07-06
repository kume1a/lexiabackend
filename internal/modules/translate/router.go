package translate

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func Router(apiCfg *shared.ApiConfig, rg *gin.RouterGroup) {
	translateGroup := rg.Group("/translate")
	{
		translateGroup.POST("", handleTranslate(apiCfg))
		translateGroup.POST("/detect", handleDetectLanguage(apiCfg))
		translateGroup.GET("/languages", handleGetSupportedLanguages(apiCfg))
	}
}
