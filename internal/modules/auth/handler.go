package auth

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func handleGetAuthStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		shared.ResOK(c, shared.OkDTO{Ok: true})
	}
}

func handleEmailSignIn(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body emailSignInDTO
		if err := c.ShouldBindJSON(&body); err != nil {
			shared.ResBadRequest(c, err.Error())
			return
		}

		tokenPayload, httpErr := SignInWithEmail(apiCfg, c.Request.Context(), SignInWithEmailArgs{
			Email:    body.Email,
			Password: body.Password,
		})
		if httpErr != nil {
			shared.ResHttpError(c, httpErr)
			return
		}

		shared.ResOK(c, tokenPayload)
	}
}

func handleEmailSignUp(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body EmailSignUpDTO
		if err := c.ShouldBindJSON(&body); err != nil {
			shared.ResBadRequest(c, err.Error())
			return
		}

		tokenPayload, httpErr := SignUpWithEmail(apiCfg, c.Request.Context(), SignUpWithEmailArgs{
			Username: body.Username,
			Email:    body.Email,
			Password: body.Password,
		})
		if httpErr != nil {
			shared.ResHttpError(c, httpErr)
			return
		}

		shared.ResOK(c, tokenPayload)
	}
}
