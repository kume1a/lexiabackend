package auth

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func Router(apiCfg *shared.ApiConfig, r *gin.Engine) *gin.Engine {
	r.GET("/auth/status", handleGetAuthStatus())

	r.PUT("/auth/emailSignIn", handleEmailSignIn(apiCfg))
	r.GET("/auth/emailSignUp", handleEmailSignUp(apiCfg))

	return r
}
