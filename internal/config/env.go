package config

import (
	"errors"
	"fmt"
	"lexia/internal/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("ENVIRONMENT")
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
	IsDevelopment         bool
	IsProduction          bool
	Port                  string
	DbConnectionString    string
	AccessTokenSecret     string
	AccessTokenExpSeconds int64
}

func ParseEnv() (*EnvVariables, error) {
	environment, err := getEnv("ENVIRONMENT")
	if err != nil {
		return nil, err
	}

	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}

	dbConnectionString, err := getEnv("DB_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}

	accessTokenSecret, err := getEnv("ACCESS_TOKEN_SECRET")
	if err != nil {
		return nil, err
	}

	accessTokenExpSeconds, err := getEnvInt("ACCESS_TOKEN_EXP_SECONDS")
	if err != nil {
		return nil, err
	}

	return &EnvVariables{
		IsDevelopment:         environment == "development",
		IsProduction:          environment == "production",
		Port:                  port,
		DbConnectionString:    dbConnectionString,
		AccessTokenSecret:     accessTokenSecret,
		AccessTokenExpSeconds: accessTokenExpSeconds,
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
