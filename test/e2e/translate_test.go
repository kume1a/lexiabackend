package e2etest

import (
	"lexia/test/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TranslateTestSuite struct {
	helpers.E2ETestSuite
	httpClient *helpers.HTTPClient
	authToken  string
}

func (suite *TranslateTestSuite) SetupTest() {
	suite.E2ETestSuite.SetupTest()
	suite.httpClient = helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())
	suite.authToken = helpers.GetTestAuthToken(suite.T(), suite.httpClient)
}

func TestTranslateTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateTestSuite))
}

func (suite *TranslateTestSuite) getAuthHeaders() map[string]string {
	return map[string]string{
		"Authorization": suite.authToken,
	}
}

func (suite *TranslateTestSuite) TestTranslateEnglishToGeorgian() {
	translateData := map[string]interface{}{
		"text":         "Hello, world!",
		"languageFrom": "ENGLISH",
		"languageTo":   "GEORGIAN",
	}

	resp := suite.httpClient.POST("/api/v1/translate", translateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "Hello, world!", response["originalText"])
	assert.Equal(suite.T(), "ENGLISH", response["languageFrom"])
	assert.Equal(suite.T(), "GEORGIAN", response["languageTo"])

	translations, ok := response["translations"].([]interface{})
	assert.True(suite.T(), ok)
	assert.Greater(suite.T(), len(translations), 0)

	// Check first translation
	firstTranslation := translations[0].(map[string]interface{})
	assert.Contains(suite.T(), firstTranslation, "text")
	assert.Contains(suite.T(), firstTranslation, "confidence")
	assert.Greater(suite.T(), firstTranslation["confidence"].(float64), 0.0)
}

func (suite *TranslateTestSuite) TestTranslateWithoutAuth() {
	translateData := map[string]interface{}{
		"text":         "Hello",
		"languageFrom": "ENGLISH",
		"languageTo":   "GEORGIAN",
	}

	resp := suite.httpClient.POST("/api/v1/translate", translateData, nil)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *TranslateTestSuite) TestTranslateSameLanguage() {
	translateData := map[string]interface{}{
		"text":         "Hello",
		"languageFrom": "ENGLISH",
		"languageTo":   "ENGLISH",
	}

	resp := suite.httpClient.POST("/api/v1/translate", translateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *TranslateTestSuite) TestTranslateSpanishToEnglish() {
	translateData := map[string]interface{}{
		"text":         "Hola, mundo!",
		"languageFrom": "SPANISH",
		"languageTo":   "ENGLISH",
	}

	resp := suite.httpClient.POST("/api/v1/translate", translateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "Hola, mundo!", response["originalText"])
	assert.Equal(suite.T(), "SPANISH", response["languageFrom"])
	assert.Equal(suite.T(), "ENGLISH", response["languageTo"])

	translations, ok := response["translations"].([]interface{})
	assert.True(suite.T(), ok)
	assert.Greater(suite.T(), len(translations), 0)
}
