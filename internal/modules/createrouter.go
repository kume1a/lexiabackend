package modules

import (
	"lexia/internal/config"
	"lexia/internal/logger"
	"lexia/internal/modules/auth"

	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthcheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func CreateWebserverRouter() error {
	envVars, err := config.ParseEnv()
	if err != nil {
		logger.Fatal("Failed to parse environment variables: ", err)
		return err
	}

	r := gin.Default()

	// healthcheck
	r.GET("/", HealthcheckHandler)

	// auth
	r.POST("/auth/signIn", auth.LoginHandler)

	if envVars.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := r.Run(); err != nil {
		logger.Fatal("Failed to start HTTP server: ", err)
		return err
	}

	return nil
}
