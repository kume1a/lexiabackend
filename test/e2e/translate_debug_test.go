package e2etest

import (
	"fmt"
	"lexia/test/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TranslateDebugTestSuite struct {
	helpers.E2ETestSuite
	httpClient *helpers.HTTPClient
	authToken  string
}

func (suite *TranslateDebugTestSuite) SetupTest() {
	suite.E2ETestSuite.SetupTest()
	suite.httpClient = helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())
	suite.authToken = helpers.GetTestAuthToken(suite.T(), suite.httpClient)
	fmt.Printf("Debug: Auth token is: %s\n", suite.authToken)
}

func TestTranslateDebugTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateDebugTestSuite))
}

func (suite *TranslateDebugTestSuite) getAuthHeaders() map[string]string {
	headers := map[string]string{
		"Authorization": suite.authToken,
	}
	fmt.Printf("Debug: Auth headers: %+v\n", headers)
	return headers
}

func (suite *TranslateDebugTestSuite) TestTranslateWithDebug() {
	translateData := map[string]interface{}{
		"text":         "Hello",
		"languageFrom": "ENGLISH",
		"languageTo":   "GEORGIAN",
	}

	fmt.Printf("Debug: Making request to /api/v1/translate\n")
	resp := suite.httpClient.POST("/api/v1/translate", translateData, suite.getAuthHeaders())
	fmt.Printf("Debug: Response status code: %d\n", resp.StatusCode)

	responseBody := resp.Body
	fmt.Printf("Debug: Response body: %s\n", string(responseBody))

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}
