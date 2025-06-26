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
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
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
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
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
