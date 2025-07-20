package shared

import (
	"errors"
	"fmt"
	"lexia/internal/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	EnvEnvironment                 = "ENVIRONMENT"
	EnvPort                        = "PORT"
	EnvDbConnectionString          = "DB_CONNECTION_STRING"
	EnvAccessTokenSecret           = "ACCESS_TOKEN_SECRET"
	EnvAccessTokenExpSeconds       = "ACCESS_TOKEN_EXP_SECONDS"
	EnvGoogleCloudProjectID        = "GOOGLE_CLOUD_PROJECT_ID"
	EnvGoogleServiceAccountKeyPath = "GOOGLE_SERVICE_ACCOUNT_KEY_PATH"
)

func LoadEnv() {
	env := os.Getenv(EnvEnvironment)
	if env == "" {
		env = "development"
	}

	if env == "development" {
		envPath := ".env." + env

		logger.Info("Loading env file: " + envPath)

		godotenv.Load(envPath)
	}
}

type EnvVariables struct {
	IsDevelopment               bool
	IsProduction                bool
	Port                        string
	DbConnectionString          string
	AccessTokenSecret           string
	AccessTokenExpSeconds       int64
	GoogleCloudProjectID        string
	GoogleServiceAccountKeyPath string
}

func ParseEnv() (*EnvVariables, error) {
	environment, err := getEnv(EnvEnvironment)
	if err != nil {
		return nil, err
	}

	port, err := getEnv(EnvPort)
	if err != nil {
		return nil, err
	}

	dbConnectionString, err := getEnv(EnvDbConnectionString)
	if err != nil {
		return nil, err
	}

	accessTokenSecret, err := getEnv(EnvAccessTokenSecret)
	if err != nil {
		return nil, err
	}

	accessTokenExpSeconds, err := getEnvInt(EnvAccessTokenExpSeconds)
	if err != nil {
		return nil, err
	}

	googleCloudProjectID, err := getEnv(EnvGoogleCloudProjectID)
	if err != nil {
		return nil, err
	}

	googleServiceAccountKeyPath := os.Getenv(EnvGoogleServiceAccountKeyPath)

	return &EnvVariables{
		IsDevelopment:               environment == "development",
		IsProduction:                environment == "production",
		Port:                        port,
		DbConnectionString:          dbConnectionString,
		AccessTokenSecret:           accessTokenSecret,
		AccessTokenExpSeconds:       accessTokenExpSeconds,
		GoogleCloudProjectID:        googleCloudProjectID,
		GoogleServiceAccountKeyPath: googleServiceAccountKeyPath,
	}, nil
}

func getEnvInt(key string) (int64, error) {
	val, err := getEnv(key)
	if err != nil {
		return 0, err
	}

	valInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return valInt, nil
}

func getEnv(key string) (string, error) {
	envVar := os.Getenv(key)
	if envVar == "" {
		msg := fmt.Sprintf("%v is not found in the env", key)

		logger.Fatal(msg)
		return "", errors.New(msg)
	}

	return envVar, nil
}
