package auth

import (
	"lexia/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": shared.ErrInvalidRequest})
		return
	}

	accessToken, err := shared.GenerateAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": shared.ErrInternal})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{AccessToken: accessToken})
}
