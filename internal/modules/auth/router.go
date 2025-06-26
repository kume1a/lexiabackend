package auth

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func Router(apiCfg *shared.ApiConfig, rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.GET("/status", handleGetAuthStatus())
		authGroup.POST("/signin", handleEmailSignIn(apiCfg))
		authGroup.POST("/signup", handleEmailSignUp(apiCfg))
	}
}
