package modules

import (
	"lexia/internal/logger"
	"lexia/internal/modules/auth"
	"lexia/internal/modules/user"
	"lexia/internal/shared"

	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthcheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func CreateWebserver(apiCfg *shared.ApiConfig) (*gin.Engine, error) {
	envVars, err := shared.ParseEnv()
	if err != nil {
		logger.Fatal("Failed to parse environment variables: ", err)
		return nil, err
	}

	r := gin.Default()

	r.GET("/", HealthcheckHandler)

	r = auth.Router(apiCfg, r)
	r = user.Router(apiCfg, r)

	if envVars.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	return r, nil
}
