package modules

import (
	"lexia/internal/logger"
	"lexia/internal/modules/auth"
	"lexia/internal/modules/folder"
	"lexia/internal/modules/translate"
	"lexia/internal/modules/user"
	"lexia/internal/modules/word"
	"lexia/internal/shared"

	"github.com/gin-gonic/gin"
)

func HealthcheckHandler(c *gin.Context) {
	shared.ResOK(c, shared.OkDTO{Ok: true})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func CreateWebserver(apiCfg *shared.ApiConfig) (*gin.Engine, error) {
	envVars, err := shared.ParseEnv()
	if err != nil {
		logger.Fatal("Failed to parse environment variables: ", err)
		return nil, err
	}

	if envVars.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", HealthcheckHandler)
	r.GET("/health", HealthcheckHandler)

	v1 := r.Group("/api/v1")
	{
		auth.Router(apiCfg, v1)

		protected := v1.Group("/")
		protected.Use(shared.AuthMW())
		{
			user.Router(apiCfg, protected)
			folder.Router(apiCfg, protected)
			word.Router(apiCfg, protected)
			translate.Router(apiCfg, protected)
		}
	}

	return r, nil
}
