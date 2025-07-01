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

type FolderTestSuite struct {
	helpers.E2ETestSuite
	httpClient *helpers.HTTPClient
	authToken  string
}

func (suite *FolderTestSuite) SetupTest() {
	suite.E2ETestSuite.SetupTest()
	suite.httpClient = helpers.NewTestHTTPClient(suite.T(), suite.GetTestServerURL())
	suite.authToken = helpers.GetTestAuthToken(suite.T(), suite.httpClient)
}

func TestFolderTestSuite(t *testing.T) {
	suite.Run(t, new(FolderTestSuite))
}

func (suite *FolderTestSuite) getAuthHeaders() map[string]string {
	return map[string]string{
		"Authorization": suite.authToken,
	}
}

func (suite *FolderTestSuite) TestCreateFolderCollection() {
	folderData := map[string]interface{}{
		"name": "My Folder Collection",
		"type": "FOLDER_COLLECTION",
	}

	resp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "My Folder Collection", response["name"])
	assert.Equal(suite.T(), "FOLDER_COLLECTION", response["type"])
	assert.Equal(suite.T(), float64(0), response["wordCount"])
	assert.Nil(suite.T(), response["languageFrom"])
	assert.Nil(suite.T(), response["languageTo"])
	assert.Nil(suite.T(), response["parentId"])
}

func (suite *FolderTestSuite) TestCreateWordCollection() {
	folderData := map[string]interface{}{
		"name":         "Georgian Words",
		"type":         "WORD_COLLECTION",
		"languageFrom": "GEORGIAN",
		"languageTo":   "ENGLISH",
	}

	resp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "Georgian Words", response["name"])
	assert.Equal(suite.T(), "WORD_COLLECTION", response["type"])
	assert.Equal(suite.T(), "GEORGIAN", response["languageFrom"])
	assert.Equal(suite.T(), "ENGLISH", response["languageTo"])
}

func (suite *FolderTestSuite) TestCreateWordCollectionWithoutLanguageFrom() {
	folderData := map[string]interface{}{
		"name": "Invalid Word Collection",
		"type": "WORD_COLLECTION",
	}

	resp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *FolderTestSuite) TestCreateFolderCollectionWithLanguages() {
	folderData := map[string]interface{}{
		"name":         "Invalid Folder Collection",
		"type":         "FOLDER_COLLECTION",
		"languageFrom": "GEORGIAN",
	}

	resp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *FolderTestSuite) TestCreateFolderWithValidation() {
	folderData := map[string]interface{}{
		"name": "",
		"type": "FOLDER_COLLECTION",
	}

	resp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *FolderTestSuite) TestCreateSubfolder() {
	parentData := map[string]interface{}{
		"name": "Parent Collection",
		"type": "FOLDER_COLLECTION",
	}

	parentResp := suite.httpClient.POST("/api/v1/folders", parentData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, parentResp.StatusCode)

	var parentResponse map[string]interface{}
	err := parentResp.ParseJSON(&parentResponse)
	assert.NoError(suite.T(), err)
	parentID := parentResponse["id"].(string)

	childData := map[string]interface{}{
		"name":     "Child Collection",
		"type":     "FOLDER_COLLECTION",
		"parentId": parentID,
	}

	childResp := suite.httpClient.POST("/api/v1/folders", childData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, childResp.StatusCode)

	var childResponse map[string]interface{}
	err = childResp.ParseJSON(&childResponse)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "Child Collection", childResponse["name"])
	assert.Equal(suite.T(), parentID, childResponse["parentId"])
}

func (suite *FolderTestSuite) TestCreateSubfolderInWordCollection() {
	parentData := map[string]interface{}{
		"name":         "Word Collection",
		"type":         "WORD_COLLECTION",
		"languageFrom": "ENGLISH",
	}

	parentResp := suite.httpClient.POST("/api/v1/folders", parentData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, parentResp.StatusCode)

	var parentResponse map[string]interface{}
	err := parentResp.ParseJSON(&parentResponse)
	assert.NoError(suite.T(), err)
	parentID := parentResponse["id"].(string)

	childData := map[string]interface{}{
		"name":     "Child Collection",
		"type":     "FOLDER_COLLECTION",
		"parentId": parentID,
	}

	childResp := suite.httpClient.POST("/api/v1/folders", childData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, childResp.StatusCode)
}

func (suite *FolderTestSuite) TestGetFolder() {
	folderData := map[string]interface{}{
		"name": "Test Folder",
		"type": "FOLDER_COLLECTION",
	}

	createResp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, createResp.StatusCode)

	var createResponse map[string]interface{}
	err := createResp.ParseJSON(&createResponse)
	assert.NoError(suite.T(), err)
	folderID := createResponse["id"].(string)

	getResp := suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s", folderID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, getResp.StatusCode)

	var getResponse map[string]interface{}
	err = getResp.ParseJSON(&getResponse)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), folderID, getResponse["id"])
	assert.Equal(suite.T(), "Test Folder", getResponse["name"])
	assert.Equal(suite.T(), "FOLDER_COLLECTION", getResponse["type"])
}

func (suite *FolderTestSuite) TestGetNonExistentFolder() {
	nonExistentID := uuid.New().String()
	resp := suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s", nonExistentID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *FolderTestSuite) TestGetFolderWithInvalidID() {
	resp := suite.httpClient.GET("/api/v1/folders/invalid-id", suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *FolderTestSuite) TestGetUserFolders() {
	folder1Data := map[string]interface{}{
		"name": "Folder 1",
		"type": "FOLDER_COLLECTION",
	}

	folder2Data := map[string]interface{}{
		"name":         "Folder 2",
		"type":         "WORD_COLLECTION",
		"languageFrom": "GEORGIAN",
	}

	suite.httpClient.POST("/api/v1/folders", folder1Data, suite.getAuthHeaders())
	suite.httpClient.POST("/api/v1/folders", folder2Data, suite.getAuthHeaders())

	resp := suite.httpClient.GET("/api/v1/folders", suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response []map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Len(suite.T(), response, 2)
}

func (suite *FolderTestSuite) TestGetRootFolders() {
	parentData := map[string]interface{}{
		"name": "Root Folder",
		"type": "FOLDER_COLLECTION",
	}

	parentResp := suite.httpClient.POST("/api/v1/folders", parentData, suite.getAuthHeaders())
	var parentResponse map[string]interface{}
	parentResp.ParseJSON(&parentResponse)
	parentID := parentResponse["id"].(string)

	childData := map[string]interface{}{
		"name":     "Child Folder",
		"type":     "FOLDER_COLLECTION",
		"parentId": parentID,
	}

	suite.httpClient.POST("/api/v1/folders", childData, suite.getAuthHeaders())

	resp := suite.httpClient.GET("/api/v1/folders/root", suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response []map[string]interface{}
	err := resp.ParseJSON(&response)
	assert.NoError(suite.T(), err)

	assert.Len(suite.T(), response, 1)
	assert.Equal(suite.T(), "Root Folder", response[0]["name"])

	subfolders := response[0]["subfolders"].([]interface{})
	assert.Len(suite.T(), subfolders, 1)
	assert.Equal(suite.T(), "Child Folder", subfolders[0].(map[string]interface{})["name"])
}

func (suite *FolderTestSuite) TestUpdateFolder() {
	folderData := map[string]interface{}{
		"name": "Original Name",
		"type": "FOLDER_COLLECTION",
	}

	createResp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	var createResponse map[string]interface{}
	createResp.ParseJSON(&createResponse)
	folderID := createResponse["id"].(string)

	updateData := map[string]interface{}{
		"name": "Updated Name",
	}

	updateResp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/folders/%s", folderID), updateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, updateResp.StatusCode)

	var updateResponse map[string]interface{}
	err := updateResp.ParseJSON(&updateResponse)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "Updated Name", updateResponse["name"])
	assert.Equal(suite.T(), folderID, updateResponse["id"])
}

func (suite *FolderTestSuite) TestUpdateFolderParent() {
	parent1Data := map[string]interface{}{
		"name": "Parent 1",
		"type": "FOLDER_COLLECTION",
	}

	parent2Data := map[string]interface{}{
		"name": "Parent 2",
		"type": "FOLDER_COLLECTION",
	}

	childData := map[string]interface{}{
		"name": "Child",
		"type": "FOLDER_COLLECTION",
	}

	parent1Resp := suite.httpClient.POST("/api/v1/folders", parent1Data, suite.getAuthHeaders())
	parent2Resp := suite.httpClient.POST("/api/v1/folders", parent2Data, suite.getAuthHeaders())
	childResp := suite.httpClient.POST("/api/v1/folders", childData, suite.getAuthHeaders())

	var parent1Response, parent2Response, childResponse map[string]interface{}
	parent1Resp.ParseJSON(&parent1Response)
	parent2Resp.ParseJSON(&parent2Response)
	childResp.ParseJSON(&childResponse)

	parent1ID := parent1Response["id"].(string)
	parent2ID := parent2Response["id"].(string)
	childID := childResponse["id"].(string)

	updateData := map[string]interface{}{
		"parentId": parent1ID,
	}

	updateResp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/folders/%s", childID), updateData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, updateResp.StatusCode)

	var updateResponse map[string]interface{}
	updateResp.ParseJSON(&updateResponse)
	assert.Equal(suite.T(), parent1ID, updateResponse["parentId"])

	updateData2 := map[string]interface{}{
		"parentId": parent2ID,
	}

	updateResp2 := suite.httpClient.PUT(fmt.Sprintf("/api/v1/folders/%s", childID), updateData2, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, updateResp2.StatusCode)

	var updateResponse2 map[string]interface{}
	updateResp2.ParseJSON(&updateResponse2)
	assert.Equal(suite.T(), parent2ID, updateResponse2["parentId"])
}

func (suite *FolderTestSuite) TestMoveFolder() {
	parent1Data := map[string]interface{}{
		"name": "Parent 1",
		"type": "FOLDER_COLLECTION",
	}

	parent2Data := map[string]interface{}{
		"name": "Parent 2",
		"type": "FOLDER_COLLECTION",
	}

	childData := map[string]interface{}{
		"name": "Child",
		"type": "FOLDER_COLLECTION",
	}

	parent1Resp := suite.httpClient.POST("/api/v1/folders", parent1Data, suite.getAuthHeaders())
	parent2Resp := suite.httpClient.POST("/api/v1/folders", parent2Data, suite.getAuthHeaders())
	childResp := suite.httpClient.POST("/api/v1/folders", childData, suite.getAuthHeaders())

	var parent1Response, parent2Response, childResponse map[string]interface{}
	parent1Resp.ParseJSON(&parent1Response)
	parent2Resp.ParseJSON(&parent2Response)
	childResp.ParseJSON(&childResponse)

	parent2ID := parent2Response["id"].(string)
	childID := childResponse["id"].(string)

	moveData := map[string]interface{}{
		"parentId": parent2ID,
	}

	moveResp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/folders/%s/move", childID), moveData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, moveResp.StatusCode)

	var moveResponse map[string]interface{}
	moveResp.ParseJSON(&moveResponse)
	assert.Equal(suite.T(), parent2ID, moveResponse["parentId"])
}

func (suite *FolderTestSuite) TestMoveFolderToRoot() {
	parentData := map[string]interface{}{
		"name": "Parent",
		"type": "FOLDER_COLLECTION",
	}

	parentResp := suite.httpClient.POST("/api/v1/folders", parentData, suite.getAuthHeaders())
	var parentResponse map[string]interface{}
	parentResp.ParseJSON(&parentResponse)
	parentID := parentResponse["id"].(string)

	childData := map[string]interface{}{
		"name":     "Child",
		"type":     "FOLDER_COLLECTION",
		"parentId": parentID,
	}

	childResp := suite.httpClient.POST("/api/v1/folders", childData, suite.getAuthHeaders())
	var childResponse map[string]interface{}
	childResp.ParseJSON(&childResponse)
	childID := childResponse["id"].(string)

	moveData := map[string]interface{}{
		"parentId": nil,
	}

	moveResp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/folders/%s/move", childID), moveData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, moveResp.StatusCode)

	var moveResponse map[string]interface{}
	moveResp.ParseJSON(&moveResponse)
	assert.Nil(suite.T(), moveResponse["parentId"])
}

func (suite *FolderTestSuite) TestDeleteEmptyFolder() {
	folderData := map[string]interface{}{
		"name": "Folder to Delete",
		"type": "FOLDER_COLLECTION",
	}

	createResp := suite.httpClient.POST("/api/v1/folders", folderData, suite.getAuthHeaders())
	var createResponse map[string]interface{}
	createResp.ParseJSON(&createResponse)
	folderID := createResponse["id"].(string)

	deleteResp := suite.httpClient.DELETE(fmt.Sprintf("/api/v1/folders/%s", folderID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusNoContent, deleteResp.StatusCode)

	getResp := suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s", folderID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusNotFound, getResp.StatusCode)
}

func (suite *FolderTestSuite) TestDeleteFolderWithSubfolders() {
	parentData := map[string]interface{}{
		"name": "Parent",
		"type": "FOLDER_COLLECTION",
	}

	parentResp := suite.httpClient.POST("/api/v1/folders", parentData, suite.getAuthHeaders())
	var parentResponse map[string]interface{}
	parentResp.ParseJSON(&parentResponse)
	parentID := parentResponse["id"].(string)

	childData := map[string]interface{}{
		"name":     "Child",
		"type":     "FOLDER_COLLECTION",
		"parentId": parentID,
	}

	suite.httpClient.POST("/api/v1/folders", childData, suite.getAuthHeaders())

	deleteResp := suite.httpClient.DELETE(fmt.Sprintf("/api/v1/folders/%s", parentID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusBadRequest, deleteResp.StatusCode)
}

func (suite *FolderTestSuite) TestDeleteNonExistentFolder() {
	nonExistentID := uuid.New().String()
	resp := suite.httpClient.DELETE(fmt.Sprintf("/api/v1/folders/%s", nonExistentID), suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func (suite *FolderTestSuite) TestCircularReferencePreventionInMove() {
	grandparentData := map[string]interface{}{
		"name": "Grandparent",
		"type": "FOLDER_COLLECTION",
	}

	grandparentResp := suite.httpClient.POST("/api/v1/folders", grandparentData, suite.getAuthHeaders())
	var grandparentResponse map[string]interface{}
	grandparentResp.ParseJSON(&grandparentResponse)
	grandparentID := grandparentResponse["id"].(string)

	parentData := map[string]interface{}{
		"name":     "Parent",
		"type":     "FOLDER_COLLECTION",
		"parentId": grandparentID,
	}

	parentResp := suite.httpClient.POST("/api/v1/folders", parentData, suite.getAuthHeaders())
	var parentResponse map[string]interface{}
	parentResp.ParseJSON(&parentResponse)
	parentID := parentResponse["id"].(string)

	childData := map[string]interface{}{
		"name":     "Child",
		"type":     "FOLDER_COLLECTION",
		"parentId": parentID,
	}

	childResp := suite.httpClient.POST("/api/v1/folders", childData, suite.getAuthHeaders())
	var childResponse map[string]interface{}
	childResp.ParseJSON(&childResponse)
	childID := childResponse["id"].(string)

	moveData := map[string]interface{}{
		"parentId": childID,
	}

	moveResp := suite.httpClient.PUT(fmt.Sprintf("/api/v1/folders/%s/move", grandparentID), moveData, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusInternalServerError, moveResp.StatusCode)
}

func (suite *FolderTestSuite) TestUnauthorizedAccess() {
	folderData := map[string]interface{}{
		"name": "Test Folder",
		"type": "FOLDER_COLLECTION",
	}

	resp := suite.httpClient.POST("/api/v1/folders", folderData)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	resp = suite.httpClient.GET("/api/v1/folders")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	resp = suite.httpClient.GET("/api/v1/folders/root")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	folderID := uuid.New().String()
	resp = suite.httpClient.GET(fmt.Sprintf("/api/v1/folders/%s", folderID))
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	resp = suite.httpClient.PUT(fmt.Sprintf("/api/v1/folders/%s", folderID), folderData)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	resp = suite.httpClient.PUT(fmt.Sprintf("/api/v1/folders/%s/move", folderID), folderData)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	resp = suite.httpClient.DELETE(fmt.Sprintf("/api/v1/folders/%s", folderID))
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *FolderTestSuite) TestCreateFolderHierarchy() {
	languagesCollection := map[string]interface{}{
		"name": "Languages",
		"type": "FOLDER_COLLECTION",
	}

	langResp := suite.httpClient.POST("/api/v1/folders", languagesCollection, suite.getAuthHeaders())
	var langResponse map[string]interface{}
	langResp.ParseJSON(&langResponse)
	langID := langResponse["id"].(string)

	georgianCollection := map[string]interface{}{
		"name":     "Georgian",
		"type":     "FOLDER_COLLECTION",
		"parentId": langID,
	}

	georgianResp := suite.httpClient.POST("/api/v1/folders", georgianCollection, suite.getAuthHeaders())
	var georgianResponse map[string]interface{}
	georgianResp.ParseJSON(&georgianResponse)
	georgianID := georgianResponse["id"].(string)

	basicWords := map[string]interface{}{
		"name":         "Basic Words",
		"type":         "WORD_COLLECTION",
		"languageFrom": "GEORGIAN",
		"languageTo":   "ENGLISH",
		"parentId":     georgianID,
	}

	basicResp := suite.httpClient.POST("/api/v1/folders", basicWords, suite.getAuthHeaders())
	assert.Equal(suite.T(), http.StatusOK, basicResp.StatusCode)

	var basicResponse map[string]interface{}
	basicResp.ParseJSON(&basicResponse)

	assert.Equal(suite.T(), "Basic Words", basicResponse["name"])
	assert.Equal(suite.T(), "WORD_COLLECTION", basicResponse["type"])
	assert.Equal(suite.T(), "GEORGIAN", basicResponse["languageFrom"])
	assert.Equal(suite.T(), "ENGLISH", basicResponse["languageTo"])
	assert.Equal(suite.T(), georgianID, basicResponse["parentId"])

	rootResp := suite.httpClient.GET("/api/v1/folders/root", suite.getAuthHeaders())
	var rootResponse []map[string]interface{}
	rootResp.ParseJSON(&rootResponse)

	assert.Len(suite.T(), rootResponse, 1)
	assert.Equal(suite.T(), "Languages", rootResponse[0]["name"])

	subfolders := rootResponse[0]["subfolders"].([]interface{})
	assert.Len(suite.T(), subfolders, 1)

	georgianFolder := subfolders[0].(map[string]interface{})
	assert.Equal(suite.T(), "Georgian", georgianFolder["name"])

	georgianSubfolders := georgianFolder["subfolders"].([]interface{})
	assert.Len(suite.T(), georgianSubfolders, 1)

	basicWordsFolder := georgianSubfolders[0].(map[string]interface{})
	assert.Equal(suite.T(), "Basic Words", basicWordsFolder["name"])
	assert.Equal(suite.T(), "WORD_COLLECTION", basicWordsFolder["type"])
}
