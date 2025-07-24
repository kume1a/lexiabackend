package e2etest

import (
	"fmt"
	"lexia/test/helpers"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WordTestSuite struct {
	helpers.E2ETestSuite
	httpClient *helpers.HTTPClient
	authToken  string
}

func (suite *WordTestSuite) SetupTest() {
	suite.E2ETestSuite.SetupTest()
	suite.httpClient = helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())
	suite.authToken = helpers.GetTestAuthToken(suite.T(), suite.httpClient)
}

func TestWordTestSuite(t *testing.T) {
	suite.Run(t, new(WordTestSuite))
}

func (suite *WordTestSuite) getAuthHeaders() map[string]string {
	return map[string]string{
		"Authorization": suite.authToken,
	}
}

func (suite *WordTestSuite) createTestFolder() string {
	folderData := map[string]interface{}{
		"name":         "Test Vocabulary Folder",
		"type":         "WORD_COLLECTION",
		"languageFrom": "ENGLISH",
		"languageTo":   "GEORGIAN",
	}

	resp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	folderID, ok := response["id"].(string)
	assert.True(suite.T(), ok)
	assert.NotEmpty(suite.T(), folderID)

	return folderID
}

func (suite *WordTestSuite) TestCreateWord() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "hello",
		"definition": "a greeting",
		"folderId":   folderID,
	}

	resp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "hello", response["text"])
	assert.Equal(suite.T(), "a greeting", response["definition"])
	assert.NotEmpty(suite.T(), response["id"])
	assert.NotEmpty(suite.T(), response["createdAt"])
	assert.NotEmpty(suite.T(), response["updatedAt"])

	folder, ok := response["folder"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), folderID, folder["id"])
	assert.Equal(suite.T(), "Test Vocabulary Folder", folder["name"])
}

func (suite *WordTestSuite) TestCreateWordWithEmptyDefinition() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "hello",
		"definition": "",
		"folderId":   folderID,
	}

	resp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "hello", response["text"])
	assert.Equal(suite.T(), "", response["definition"])
}

func (suite *WordTestSuite) TestCreateWordValidationErrors() {
	folderID := suite.createTestFolder()

	testCases := []struct {
		name           string
		wordData       map[string]interface{}
		expectedStatus int
		expectedMsg    string
	}{
		{
			name: "missing text",
			wordData: map[string]interface{}{
				"definition": "a greeting",
				"folderId":   folderID,
			},
			expectedStatus: 500, // Database constraint violation
			expectedMsg:    "",
		},
		{
			name: "empty text",
			wordData: map[string]interface{}{
				"text":       "",
				"definition": "a greeting",
				"folderId":   folderID,
			},
			expectedStatus: 500, // Database constraint violation
			expectedMsg:    "",
		},
		{
			name: "missing folderId",
			wordData: map[string]interface{}{
				"text":       "hello",
				"definition": "a greeting",
			},
			expectedStatus: 500, // Database constraint violation (null UUID)
			expectedMsg:    "",
		},
		{
			name: "invalid folderId",
			wordData: map[string]interface{}{
				"text":       "hello",
				"definition": "a greeting",
				"folderId":   "invalid-uuid",
			},
			expectedStatus: 400,
			expectedMsg:    "Validation failed",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			resp := suite.httpClient.POST("/api/v1/words", tc.wordData, suite.getAuthHeaders())
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == 400 {
				var response map[string]interface{}
				err := resp.ParseJSON(&response)
				assert.NoError(t, err)
				assert.Contains(t, response["message"], tc.expectedMsg)
			}
		})
	}
}

func (suite *WordTestSuite) TestCreateWordWithNonExistentFolder() {
	nonExistentFolderID := uuid.New().String()

	wordData := map[string]interface{}{
		"text":       "hello",
		"definition": "a greeting",
		"folderId":   nonExistentFolderID,
	}

	resp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func (suite *WordTestSuite) TestCreateWordUnauthorized() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "hello",
		"definition": "a greeting",
		"folderId":   folderID,
	}

	resp := suite.httpClient.POST("/api/v1/words", wordData)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *WordTestSuite) TestGetWord() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "world",
		"definition": "the earth",
		"folderId":   folderID,
	}

	createResp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createResponse map[string]interface{}
	err := createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)

	wordID := createResponse["id"].(string)

	resp := suite.httpClient.GET(fmt.Sprintf("/api/v1/words/%s", wordID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), wordID, response["id"])
	assert.Equal(suite.T(), "world", response["text"])
	assert.Equal(suite.T(), "the earth", response["definition"])

	folder, ok := response["folder"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), folderID, folder["id"])
	assert.Equal(suite.T(), "Test Vocabulary Folder", folder["name"])
}

func (suite *WordTestSuite) TestGetWordNotFound() {
	nonExistentWordID := uuid.New().String()

	resp := suite.httpClient.GET(fmt.Sprintf("/api/v1/words/%s", nonExistentWordID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Word not found", response["error"])
}

func (suite *WordTestSuite) TestGetWordInvalidID() {
	resp := suite.httpClient.GET("/api/v1/words/invalid-uuid", suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid word ID", response["error"])
}

func (suite *WordTestSuite) TestGetWordsByFolder() {
	folderID := suite.createTestFolder()

	words := []map[string]interface{}{
		{
			"text":       "apple",
			"definition": "a red fruit",
			"folderId":   folderID,
		},
		{
			"text":       "banana",
			"definition": "a yellow fruit",
			"folderId":   folderID,
		},
		{
			"text":       "cherry",
			"definition": "a small red fruit",
			"folderId":   folderID,
		},
	}

	for _, wordData := range words {
		resp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	}

	resp := suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s/words", folderID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response []map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Len(suite.T(), response, 3)

	wordTexts := make([]string, len(response))
	for i, word := range response {
		wordTexts[i] = word["text"].(string)
		assert.NotEmpty(suite.T(), word["id"])
		assert.NotEmpty(suite.T(), word["createdAt"])
		assert.NotEmpty(suite.T(), word["updatedAt"])
		assert.NotEmpty(suite.T(), word["definition"])
		assert.Equal(suite.T(), folderID, word["folderId"])
	}

	assert.Contains(suite.T(), wordTexts, "apple")
	assert.Contains(suite.T(), wordTexts, "banana")
	assert.Contains(suite.T(), wordTexts, "cherry")
}

func (suite *WordTestSuite) TestGetWordsByFolderUnauthorized() {
	folderID := suite.createTestFolder()

	resp := suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s/words", folderID))
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *WordTestSuite) TestGetWordsByFolderInvalidID() {
	resp := suite.httpClient.GET("/api/v1/folders/invalid-uuid/words", suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid folder ID", response["error"])
}

func (suite *WordTestSuite) TestGetWordsByNonExistentFolder() {
	nonExistentFolderID := uuid.New().String()

	resp := suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s/words", nonExistentFolderID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func (suite *WordTestSuite) TestUpdateWord() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "old text",
		"definition": "old definition",
		"folderId":   folderID,
	}

	createResp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createResponse map[string]interface{}
	err := createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)

	wordID := createResponse["id"].(string)

	updateData := map[string]interface{}{
		"text":       "new text",
		"definition": "new definition",
	}

	resp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/words/%s", wordID), updateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), wordID, response["id"])
	assert.Equal(suite.T(), "new text", response["text"])
	assert.Equal(suite.T(), "new definition", response["definition"])

	folder, ok := response["folder"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), folderID, folder["id"])
}

func (suite *WordTestSuite) TestUpdateWordPartial() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "original text",
		"definition": "original definition",
		"folderId":   folderID,
	}

	createResp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createResponse map[string]interface{}
	err := createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)

	wordID := createResponse["id"].(string)

	updateData := map[string]interface{}{
		"text": "updated text only",
	}

	resp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/words/%s", wordID), updateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "updated text only", response["text"])
	assert.Equal(suite.T(), "original definition", response["definition"])

	updateDataDefinition := map[string]interface{}{
		"definition": "updated definition only",
	}

	resp = suite.httpClient.PUT(fmt.Sprintf("/api/v1/words/%s", wordID), updateDataDefinition, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	err = resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "updated text only", response["text"])
	assert.Equal(suite.T(), "updated definition only", response["definition"])
}

func (suite *WordTestSuite) TestUpdateWordValidationErrors() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "original text",
		"definition": "original definition",
		"folderId":   folderID,
	}

	createResp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createResponse map[string]interface{}
	err := createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)

	wordID := createResponse["id"].(string)

	testCases := []struct {
		name           string
		updateData     map[string]interface{}
		expectedStatus int
		expectedMsg    string
	}{
		{
			name: "empty_text",
			updateData: map[string]interface{}{
				"text": "",
			},
			expectedStatus: 500,
			expectedMsg:    "INTERNAL",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			resp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/words/%s", wordID), tc.updateData, suite.getAuthHeaders())
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			err := resp.ParseJSON(&response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], tc.expectedMsg)
		})
	}
}

func (suite *WordTestSuite) TestUpdateWordUnauthorized() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "test word",
		"definition": "test definition",
		"folderId":   folderID,
	}

	createResp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createResponse map[string]interface{}
	err := createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)

	wordID := createResponse["id"].(string)

	updateData := map[string]interface{}{
		"text": "unauthorized update",
	}

	resp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/words/%s", wordID), updateData)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *WordTestSuite) TestUpdateWordNotFound() {
	nonExistentWordID := uuid.New().String()

	updateData := map[string]interface{}{
		"text": "new text",
	}

	resp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/words/%s", nonExistentWordID), updateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func (suite *WordTestSuite) TestUpdateWordInvalidID() {
	updateData := map[string]interface{}{
		"text": "new text",
	}

	resp := suite.httpClient.PUT("/api/v1/words/invalid-uuid", updateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid word ID", response["error"])
}

func (suite *WordTestSuite) TestDeleteWord() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "word to delete",
		"definition": "this word will be deleted",
		"folderId":   folderID,
	}

	createResp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createResponse map[string]interface{}
	err := createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)

	wordID := createResponse["id"].(string)

	resp := suite.httpClient.DELETE(fmt.Sprintf("/api/v1/words/%s", wordID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	getResp := suite.httpClient.GET(fmt.Sprintf("/api/v1/words/%s", wordID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusNotFound, getResp.StatusCode)
}

func (suite *WordTestSuite) TestDeleteWordUnauthorized() {
	folderID := suite.createTestFolder()

	wordData := map[string]interface{}{
		"text":       "test word",
		"definition": "test definition",
		"folderId":   folderID,
	}

	createResp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createResponse map[string]interface{}
	err := createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)

	wordID := createResponse["id"].(string)

	resp := suite.httpClient.DELETE(fmt.Sprintf("/api/v1/words/%s", wordID))
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *WordTestSuite) TestDeleteWordNotFound() {
	nonExistentWordID := uuid.New().String()

	resp := suite.httpClient.DELETE(fmt.Sprintf("/api/v1/words/%s", nonExistentWordID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func (suite *WordTestSuite) TestDeleteWordInvalidID() {
	resp := suite.httpClient.DELETE("/api/v1/words/invalid-uuid", suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid word ID", response["error"])
}

func (suite *WordTestSuite) TestWordOperationsWithDifferentUsers() {
	user1Token := suite.authToken

	// Create a second user with different credentials
	user2HttpClient := helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())

	signupData := map[string]string{
		"email":    "test2@example.com",
		"password": "password123",
		"username": "testuser2",
	}

	signupResp := user2HttpClient.POST("/api/v1/auth/signup", signupData)
	assert.Equal(suite.T(), http.StatusOK, signupResp.StatusCode)

	signinData := map[string]string{
		"email":    "test2@example.com",
		"password": "password123",
	}

	signinResp := user2HttpClient.POST("/api/v1/auth/signin", signinData)
	assert.Equal(suite.T(), http.StatusOK, signinResp.StatusCode)

	var signinResponse map[string]interface{}
	signinErr := signinResp.ParseJSON(&signinResponse)
	assert.NoError(suite.T(), signinErr)
	user2Token := signinResponse["accessToken"].(string)

	user1Headers := map[string]string{"Authorization": user1Token}
	user2Headers := map[string]string{"Authorization": user2Token}

	folderData := map[string]interface{}{
		"name":         "User1 Folder",
		"type":         "WORD_COLLECTION",
		"languageFrom": "ENGLISH",
		"languageTo":   "GEORGIAN",
	}

	folderResp := suite.httpClient.POST("/api/v1/folders", folderData, user1Headers)
	assert.Equal(suite.T(), http.StatusOK, folderResp.StatusCode)

	var folderResponse map[string]interface{}
	err := folderResp.ParseJSON(&folderResponse)
	assert.NoError(suite.T(), err)
	folderID := folderResponse["id"].(string)

	wordData := map[string]interface{}{
		"text":       "private word",
		"definition": "only user1 can access",
		"folderId":   folderID,
	}

	createResp := suite.httpClient.POST("/api/v1/words", wordData, user1Headers)
	assert.Equal(suite.T(), http.StatusCreated, createResp.StatusCode)

	var createResponse map[string]interface{}
	err = createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)
	wordID := createResponse["id"].(string)

	getResp := user2HttpClient.GET(fmt.Sprintf("/api/v1/folders/%s/words", folderID), user2Headers)
	assert.Equal(suite.T(), http.StatusInternalServerError, getResp.StatusCode)

	updateData := map[string]interface{}{
		"text": "hacked word",
	}
	updateResp := user2HttpClient.PUT(fmt.Sprintf("/api/v1/words/%s", wordID), updateData, user2Headers)
	assert.Equal(suite.T(), http.StatusInternalServerError, updateResp.StatusCode)

	deleteResp := user2HttpClient.DELETE(fmt.Sprintf("/api/v1/words/%s", wordID), user2Headers)
	assert.Equal(suite.T(), http.StatusInternalServerError, deleteResp.StatusCode)

	verifyResp := suite.httpClient.GET(fmt.Sprintf("/api/v1/words/%s", wordID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, verifyResp.StatusCode)
}

func (suite *WordTestSuite) TestBulkWordOperations() {
	folderID := suite.createTestFolder()

	words := []map[string]interface{}{
		{"text": "word1", "definition": "definition1", "folderId": folderID},
		{"text": "word2", "definition": "definition2", "folderId": folderID},
		{"text": "word3", "definition": "definition3", "folderId": folderID},
		{"text": "word4", "definition": "definition4", "folderId": folderID},
		{"text": "word5", "definition": "definition5", "folderId": folderID},
	}

	createdWordIDs := make([]string, len(words))

	for i, wordData := range words {
		resp := suite.httpClient.POST("/api/v1/words", wordData, suite.getAuthHeaders())
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err := resp.ParseJSON(&response)
		assert.NoError(suite.T(), err)
		createdWordIDs[i] = response["id"].(string)
	}

	resp := suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s/words", folderID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var wordsResponse []map[string]interface{}
	err := resp.ParseJSON(&wordsResponse)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), wordsResponse, 5)

	for _, wordID := range createdWordIDs {
		deleteResp := suite.httpClient.DELETE(fmt.Sprintf("/api/v1/words/%s", wordID), suite.getAuthHeaders())
		assert.Equal(suite.T(), http.StatusNoContent, deleteResp.StatusCode)
	}

	finalResp := suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s/words", folderID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, finalResp.StatusCode)

	var finalWordsResponse []map[string]interface{}
	err = finalResp.ParseJSON(&finalWordsResponse)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), finalWordsResponse, 0)
}

func (suite *WordTestSuite) TestCheckWordDuplicate() {
	folderID := suite.createTestFolder()

	// Create a word
	wordPayload := map[string]interface{}{
		"text":       "hello",
		"definition": "a greeting",
		"folderId":   folderID,
	}

	wordResp := suite.httpClient.POST("/api/v1/words", wordPayload, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, wordResp.StatusCode)

	suite.T().Run("should find duplicate word", func(t *testing.T) {
		// Check for duplicate
		resp := suite.httpClient.GET("/api/v1/words/check-duplicate?text=hello", suite.getAuthHeaders())
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := resp.ParseJSON(&result)
		assert.NoError(t, err)

		assert.True(t, result["isDuplicate"].(bool))
		assert.NotNil(t, result["word"])

		word := result["word"].(map[string]interface{})
		assert.Equal(t, "hello", word["text"])
		assert.Equal(t, "a greeting", word["definition"])
		assert.NotEmpty(t, word["folderPath"])

		folderPath := word["folderPath"].([]interface{})
		assert.Len(t, folderPath, 1)

		folder := folderPath[0].(map[string]interface{})
		assert.Equal(t, "Test Vocabulary Folder", folder["name"])
	})

	suite.T().Run("should not find duplicate for non-existent word", func(t *testing.T) {
		// Check for non-existent word
		resp := suite.httpClient.GET("/api/v1/words/check-duplicate?text=nonexistent", suite.getAuthHeaders())
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := resp.ParseJSON(&result)
		assert.NoError(t, err)

		assert.False(t, result["isDuplicate"].(bool))
		assert.Nil(t, result["word"])
	})

	suite.T().Run("should return bad request when text parameter is missing", func(t *testing.T) {
		resp := suite.httpClient.GET("/api/v1/words/check-duplicate", suite.getAuthHeaders())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	suite.T().Run("should return unauthorized when no auth token", func(t *testing.T) {
		resp := suite.httpClient.GET("/api/v1/words/check-duplicate?text=hello")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func (suite *WordTestSuite) TestCheckWordDuplicateWithNestedFolders() {
	// Create parent folder
	parentFolderData := map[string]interface{}{
		"name": "Parent Folder",
		"type": "FOLDER_COLLECTION",
	}

	parentResp := suite.httpClient.POST("/api/v1/folders", parentFolderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, parentResp.StatusCode)

	var parentFolderResult map[string]interface{}
	err := parentResp.ParseJSON(&parentFolderResult)
	assert.NoError(suite.T(), err)
	parentFolderID := parentFolderResult["id"].(string)

	// Create child folder
	childFolderData := map[string]interface{}{
		"name":         "Child Folder",
		"type":         "WORD_COLLECTION",
		"languageFrom": "ENGLISH",
		"languageTo":   "SPANISH",
		"parentId":     parentFolderID,
	}

	childResp := suite.httpClient.POST("/api/v1/folders", childFolderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, childResp.StatusCode)

	var childFolderResult map[string]interface{}
	err = childResp.ParseJSON(&childFolderResult)
	assert.NoError(suite.T(), err)
	childFolderID := childFolderResult["id"].(string)

	// Create a word in the child folder
	wordPayload := map[string]interface{}{
		"text":       "nested",
		"definition": "word in nested folder",
		"folderId":   childFolderID,
	}

	wordResp := suite.httpClient.POST("/api/v1/words", wordPayload, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, wordResp.StatusCode)

	// Check for duplicate
	resp := suite.httpClient.GET("/api/v1/words/check-duplicate?text=nested", suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = resp.ParseJSON(&result)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), result["isDuplicate"].(bool))
	assert.NotNil(suite.T(), result["word"])

	word := result["word"].(map[string]interface{})
	assert.Equal(suite.T(), "nested", word["text"])
	assert.Equal(suite.T(), "word in nested folder", word["definition"])

	folderPath := word["folderPath"].([]interface{})
	assert.Len(suite.T(), folderPath, 2) // Should show parent -> child path

	// First folder should be parent
	parentFolder := folderPath[0].(map[string]interface{})
	assert.Equal(suite.T(), "Parent Folder", parentFolder["name"])
	assert.Equal(suite.T(), parentFolderID, parentFolder["id"])

	// Second folder should be child
	childFolder := folderPath[1].(map[string]interface{})
	assert.Equal(suite.T(), "Child Folder", childFolder["name"])
	assert.Equal(suite.T(), childFolderID, childFolder["id"])
}

func (suite *WordTestSuite) TestCheckWordDuplicateUserIsolation() {
	// Create a word with the first user
	folderID := suite.createTestFolder()
	wordPayload := map[string]interface{}{
		"text":       "isolation_test",
		"definition": "test word for user isolation",
		"folderId":   folderID,
	}

	wordResp := suite.httpClient.POST("/api/v1/words", wordPayload, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusCreated, wordResp.StatusCode)

	// Create a second user
	user2HttpClient := helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())

	signupData := map[string]string{
		"email":    "user2@example.com",
		"password": "password123",
		"username": "testuser2",
	}

	signupResp := user2HttpClient.POST("/api/v1/auth/signup", signupData)
	assert.Equal(suite.T(), http.StatusOK, signupResp.StatusCode)

	signinData := map[string]string{
		"email":    "user2@example.com",
		"password": "password123",
	}

	signinResp := user2HttpClient.POST("/api/v1/auth/signin", signinData)
	assert.Equal(suite.T(), http.StatusOK, signinResp.StatusCode)

	var signinResponse map[string]interface{}
	err := signinResp.ParseJSON(&signinResponse)
	assert.NoError(suite.T(), err)
	user2Token := signinResponse["accessToken"].(string)

	user2Headers := map[string]string{"Authorization": user2Token}

	// User 2 checks for the same word text - should not find it
	resp := user2HttpClient.GET("/api/v1/words/check-duplicate?text=isolation_test", user2Headers)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = resp.ParseJSON(&result)
	assert.NoError(suite.T(), err)

	assert.False(suite.T(), result["isDuplicate"].(bool))
	assert.Nil(suite.T(), result["word"])

	// User 1 should still find the word
	resp = suite.httpClient.GET("/api/v1/words/check-duplicate?text=isolation_test", suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	err = resp.ParseJSON(&result)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), result["isDuplicate"].(bool))
	assert.NotNil(suite.T(), result["word"])
}
