package user

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func Router(apiCfg *shared.ApiConfig, rg *gin.RouterGroup) {
	userGroup := rg.Group("/user")
	{
		userGroup.GET("/auth", handleGetAuthUser(apiCfg))
		userGroup.PUT("/", handleUpdateUser(apiCfg))
	}
}
