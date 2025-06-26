package shared

import (
	"github.com/gin-gonic/gin"
)

func ValidationMiddleware[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body T
		if validationErr := BindAndValidate(c, &body); validationErr != nil {
			ResValidationError(c, validationErr)
			c.Abort()
			return
		}

		c.Set("validatedBody", body)
		c.Next()
	}
}

func GetValidatedBody[T any](c *gin.Context) (T, bool) {
	body, exists := c.Get("validatedBody")
	if !exists {
		var zero T
		return zero, false
	}

	typedBody, ok := body.(T)
	return typedBody, ok
}
