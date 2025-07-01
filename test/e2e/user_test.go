package e2etest

import (
	"lexia/test/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserTestSuite struct {
	helpers.E2ETestSuite
	httpClient *helpers.HTTPClient
}

func (suite *UserTestSuite) SetupTest() {
	suite.E2ETestSuite.SetupTest()
	suite.httpClient = helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())
}

func (suite *UserTestSuite) TestGetAuthUserWithoutToken() {
	resp := suite.httpClient.GET("/api/v1/user/auth")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *UserTestSuite) TestGetAuthUserWithValidToken() {
	token := helpers.GetTestAuthToken(suite.T(), suite.httpClient)

	headers := map[string]string{
		"Authorization": token,
	}
	resp := suite.httpClient.GET("/api/v1/user/auth", headers)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]any
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Contains(suite.T(), response, "email")
	assert.Contains(suite.T(), response, "username")
	assert.Equal(suite.T(), "test@example.com", response["email"])
	assert.Equal(suite.T(), "testuser", response["username"])
}

func (suite *UserTestSuite) TestGetAuthUserWithInvalidToken() {
	headers := map[string]string{
		"Authorization": "Bearer invalid-token",
	}
	resp := suite.httpClient.GET("/api/v1/user/auth", headers)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *UserTestSuite) TestUpdateUserWithoutToken() {
	updateData := map[string]string{
		"username": "newusername",
	}

	resp := suite.httpClient.PUT("/api/v1/user/auth", updateData)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *UserTestSuite) TestUpdateUserWithValidToken() {
	token := helpers.GetTestAuthToken(suite.T(), suite.httpClient)

	updateData := map[string]string{
		"username": "newusername",
	}

	headers := map[string]string{
		"Authorization": token,
	}
	resp := suite.httpClient.PUT("/api/v1/user/auth", updateData, headers)

	assert.True(suite.T(), resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent)

	users, err := suite.GetDBClient().User.Query().All(suite.GetContext())
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), "newusername", users[0].Username)
}

func (suite *UserTestSuite) TestUpdateUserWithInvalidToken() {
	updateData := map[string]string{
		"username": "newusername",
	}

	headers := map[string]string{
		"Authorization": "Bearer invalid-token",
	}
	resp := suite.httpClient.PUT("/api/v1/user/auth", updateData, headers)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func TestUserSuite(t *testing.T) {
	helpers.RunE2ETestSuite(t, new(UserTestSuite))
}
