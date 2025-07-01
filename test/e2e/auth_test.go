package e2etest

import (
	"lexia/test/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AuthTestSuite struct {
	helpers.E2ETestSuite
	httpClient *helpers.HTTPClient
}

func (suite *AuthTestSuite) SetupTest() {
	suite.E2ETestSuite.SetupTest()
	suite.httpClient = helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())
}

func (suite *AuthTestSuite) TestAuthStatus() {
	resp := suite.httpClient.GET("/api/v1/auth/status")
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]any
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)
}

func (suite *AuthTestSuite) TestSignUpSuccess() {
	signupData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"username": "testuser",
	}

	resp := suite.httpClient.POST("/api/v1/auth/signup", signupData)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]any
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	users, err := suite.GetDBClient().User.Query().All(suite.GetContext())
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), "test@example.com", users[0].Email)
	assert.Equal(suite.T(), "testuser", users[0].Username)
}

func (suite *AuthTestSuite) TestSignUpDuplicateEmail() {
	signupData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"username": "testuser",
	}

	resp := suite.httpClient.POST("/api/v1/auth/signup", signupData)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	resp = suite.httpClient.POST("/api/v1/auth/signup", signupData)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *AuthTestSuite) TestSignUpInvalidData() {
	testCases := []struct {
		name         string
		data         map[string]string
		expectedCode int
	}{
		{
			name:         "missing email",
			data:         map[string]string{"password": "password123", "username": "testuser"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing password",
			data:         map[string]string{"email": "test@example.com", "username": "testuser"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing username",
			data:         map[string]string{"email": "test@example.com", "password": "password123"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid email",
			data:         map[string]string{"email": "invalid-email", "password": "password123", "username": "testuser"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "short password",
			data:         map[string]string{"email": "test@example.com", "password": "123", "username": "testuser"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			resp := suite.httpClient.POST("/api/v1/auth/signup", tc.data)
			assert.Equal(t, tc.expectedCode, resp.StatusCode)
		})
	}
}

func (suite *AuthTestSuite) TestSignInSuccess() {
	signupData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"username": "testuser",
	}

	resp := suite.httpClient.POST("/api/v1/auth/signup", signupData)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	signinData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	resp = suite.httpClient.POST("/api/v1/auth/signin", signinData)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]any
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Contains(suite.T(), response, "accessToken")
	assert.NotEmpty(suite.T(), response["accessToken"])
}

func (suite *AuthTestSuite) TestSignInInvalidCredentials() {
	signinData := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}

	resp := suite.httpClient.POST("/api/v1/auth/signin", signinData)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	signupData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"username": "testuser",
	}

	resp = suite.httpClient.POST("/api/v1/auth/signup", signupData)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	signinData = map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}

	resp = suite.httpClient.POST("/api/v1/auth/signin", signinData)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *AuthTestSuite) TestCORSHeaders() {
	resp := suite.httpClient.Do(helpers.Request{
		Method: "OPTIONS",
		Path:   "/api/v1/auth/status",
	})

	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)
	assert.Equal(suite.T(), "*", resp.Headers.Get("Access-Control-Allow-Origin"))
	assert.Contains(suite.T(), resp.Headers.Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(suite.T(), resp.Headers.Get("Access-Control-Allow-Methods"), "GET")
}

func TestAuthSuite(t *testing.T) {
	helpers.RunE2ETestSuite(t, new(AuthTestSuite))
}
