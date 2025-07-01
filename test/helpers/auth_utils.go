package helpers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetTestAuthToken(t *testing.T, httpClient *HTTPClient) string {
	signupData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"username": "testuser",
	}

	resp := httpClient.POST("/api/v1/auth/signup", signupData)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	signinData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	resp = httpClient.POST("/api/v1/auth/signin", signinData)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]any
	err := resp.ParseJSON(&response)
	assert.NoError(t, err)

	token, ok := response["accessToken"].(string)
	assert.True(t, ok, "Token should be a string")
	assert.NotEmpty(t, token)

	return "Bearer " + token
}
