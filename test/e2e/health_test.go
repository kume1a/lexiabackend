package e2etest

import (
	"lexia/test/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HealthTestSuite struct {
	helpers.E2ETestSuite
	httpClient *helpers.HTTPClient
}

func (suite *HealthTestSuite) SetupTest() {
	suite.E2ETestSuite.SetupTest()
	suite.httpClient = helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())
}

func (suite *HealthTestSuite) TestHealthCheck() {
	resp := suite.httpClient.GET("/")
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]any
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response["ok"].(bool))

	resp = suite.httpClient.GET("/health")
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	err = resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response["ok"].(bool))
}

func TestHealthSuite(t *testing.T) {
	helpers.RunE2ETestSuite(t, new(HealthTestSuite))
}
