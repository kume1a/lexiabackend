package user

import (
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func handleUpdateUser(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		var body updateUserDTO
		if validationErr := shared.BindAndValidate(c, &body); validationErr != nil {
			shared.ResValidationError(c, validationErr)
			return
		}

		user, err := UpdateUserByID(
			c.Request.Context(), apiCfg.DB,
			UpdateUserByIDArgs{
				Username: body.Username,
				UserID:   authPayload.UserID,
			},
		)

		if err != nil {
			shared.ResInternalServerErrorDef(c)
			return
		}

		shared.ResOK(c, UserEntityToDto(user))
	}
}

func handleGetAuthUser(apiCfg *shared.ApiConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, err := shared.GetAuthPayload(c)
		if err != nil {
			shared.ResUnauthorized(c, err.Error())
			return
		}

		user, err := GetUserByID(c.Request.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResNotFound(c, shared.ErrUserNotFound)
			return
		}

		shared.ResOK(c, UserEntityToDto(user))
	}
}
