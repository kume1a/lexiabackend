package shared

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := GetAccessTokenFromRequest(c)
		if err != nil {
			ResUnauthorized(c, err.Error())
			return
		}

		if _, err := VerifyAccessToken(accessToken); err != nil {
			ResUnauthorized(c, ErrInvalidToken)
			return
		}

		c.Next()
	}
}

func GetAuthPayload(c *gin.Context) (*TokenClaims, error) {
	accessToken, err := GetAccessTokenFromRequest(c)
	if err != nil {
		return nil, err
	}

	return VerifyAccessToken(accessToken)
}

func GetAccessTokenFromRequest(c *gin.Context) (string, error) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		return "", errors.New(ErrMissingToken)
	}

	accessToken := strings.Replace(authorizationHeader, "Bearer ", "", 1)
	if accessToken == "" {
		return "", errors.New(ErrInvalidToken)
	}

	return accessToken, nil
}
