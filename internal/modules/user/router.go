package user

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func Router(apiCfg *shared.ApiConfig, router *gin.Engine) *gin.Engine {
	router.PUT("/user", handleUpdateUser(apiCfg))
	router.GET("/user/auth", handleGetAuthUser(apiCfg))

	return router
}
