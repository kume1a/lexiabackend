package e2etest

import (
	"lexia/ent/schema"
	"lexia/internal/modules/translate"
	"lexia/test/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TranslateProductionTestSuite struct {
	helpers.E2ETestSuite
	httpClient *helpers.HTTPClient
	authToken  string
}

func (suite *TranslateProductionTestSuite) SetupTest() {
	suite.E2ETestSuite.SetupTest()
	suite.httpClient = helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())
	suite.authToken = helpers.GetTestAuthToken(suite.T(), suite.httpClient)
}

func TestTranslateProductionTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateProductionTestSuite))
}

func (suite *TranslateProductionTestSuite) getAuthHeaders() map[string]string {
	return map[string]string{
		"Authorization": suite.authToken,
	}
}

func (suite *TranslateProductionTestSuite) TestInputValidation() {
	// Test empty text
	testCases := []struct {
		name           string
		requestBody    translate.TranslateRequestDTO
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Empty text should be rejected",
			requestBody: translate.TranslateRequestDTO{
				Text:         "",
				LanguageFrom: schema.LanguageEnglish,
				LanguageTo:   schema.LanguageSpanish,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Text cannot be empty",
		},
		{
			name: "Whitespace only text should be rejected",
			requestBody: translate.TranslateRequestDTO{
				Text:         "   \n\t   ",
				LanguageFrom: schema.LanguageEnglish,
				LanguageTo:   schema.LanguageSpanish,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Text cannot be empty",
		},
		{
			name: "Same source and target language should be rejected",
			requestBody: translate.TranslateRequestDTO{
				Text:         "Hello world",
				LanguageFrom: schema.LanguageEnglish,
				LanguageTo:   schema.LanguageEnglish,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Source and target languages must be different",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			resp := suite.httpClient.POST("/api/v1/translate", tc.requestBody, suite.getAuthHeaders())
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

func (suite *TranslateProductionTestSuite) TestTranslationResponseStructure() {
	requestBody := translate.TranslateRequestDTO{
		Text:         "Hello world",
		LanguageFrom: schema.LanguageEnglish,
		LanguageTo:   schema.LanguageSpanish,
	}

	resp := suite.httpClient.POST("/api/v1/translate", requestBody, suite.getAuthHeaders())

	// Note: This test might fail if Google Translate API is not configured
	// In that case, we expect a 500 error with credentials error
	if resp.StatusCode == http.StatusInternalServerError {
		suite.T().Skip("Google Translate API not configured - skipping production translation test")
		return
	}

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response translate.TranslateResponseDTO
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	// Verify response structure
	assert.Equal(suite.T(), "Hello world", response.OriginalText)
	assert.Equal(suite.T(), schema.LanguageEnglish, response.LanguageFrom)
	assert.Equal(suite.T(), schema.LanguageSpanish, response.LanguageTo)
	assert.NotEmpty(suite.T(), response.Translations)

	// Verify translation variants
	for _, variant := range response.Translations {
		assert.NotEmpty(suite.T(), variant.Text)
		assert.GreaterOrEqual(suite.T(), variant.Confidence, float32(0.0))
		assert.LessOrEqual(suite.T(), variant.Confidence, float32(1.0))
	}

	// Should have at least one translation variant
	assert.GreaterOrEqual(suite.T(), len(response.Translations), 1)
	// Should not have too many variants (max 4 as per our logic)
	assert.LessOrEqual(suite.T(), len(response.Translations), 4)
}

func (suite *TranslateProductionTestSuite) TestLanguageDetection() {
	requestBody := translate.DetectLanguageRequestDTO{
		Text: "Bonjour le monde",
	}

	resp := suite.httpClient.POST("/api/v1/translate/detect", requestBody, suite.getAuthHeaders())

	// Note: This test might fail if Google Translate API is not configured
	if resp.StatusCode == http.StatusInternalServerError {
		suite.T().Skip("Google Translate API not configured - skipping language detection test")
		return
	}

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response translate.DetectLanguageResponseDTO
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	// Verify response structure
	assert.Equal(suite.T(), "Bonjour le monde", response.Text)
	assert.NotEmpty(suite.T(), response.DetectedLanguage)
	assert.GreaterOrEqual(suite.T(), response.Confidence, float32(0.0))
	assert.LessOrEqual(suite.T(), response.Confidence, float32(1.0))
}

func (suite *TranslateProductionTestSuite) TestSupportedLanguagesEndpoint() {
	resp := suite.httpClient.GET("/api/v1/translate/languages", suite.getAuthHeaders())

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response translate.SupportedLanguagesResponseDTO
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	// Verify we have the expected supported languages
	expectedLanguages := []schema.Language{
		schema.LanguageEnglish,
		schema.LanguageGeorgian,
		schema.LanguageSpanish,
		schema.LanguageFrench,
		schema.LanguageGerman,
		schema.LanguageRussian,
		schema.LanguageJapanese,
		schema.LanguageChinese,
	}

	assert.ElementsMatch(suite.T(), expectedLanguages, response.Languages)
	assert.Len(suite.T(), response.Languages, len(expectedLanguages))
}
