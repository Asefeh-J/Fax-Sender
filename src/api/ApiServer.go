package api

import (
	"encoding/json"
	"errors"
	"faxsender/src/utilities"
	"faxsender/src/utilities/logger"
	"io"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

// InitRouters initializes API routes on the provided Gin router.
// It defines paths for various API endpoints and assigns routes for each endpoint.
//
// Parameters:
//   - router: A pointer to the Gin router.
func InitRouters(router *gin.Engine) {

	authtenticationPath := path.Join(utilities.API_PATHS, API_UI_AUTHENTICATION)
	saveSettings := path.Join(utilities.API_PATHS, API_UI_SAVE_SETTINGS)
	loadSettings := path.Join(utilities.API_PATHS, API_UI_LOAD_SETTINGS)
	loadAccountInfo := path.Join(utilities.API_PATHS, API_UI_LOAD_ACCOUNT_INFO)
	logout := path.Join(utilities.API_PATHS, API_UI_LOGOUT)
	getLastFaxes := path.Join(utilities.API_PATHS, API_UI_GET_LAST_FAXES)
	loadAllAccounts := path.Join(utilities.API_PATHS, API_UI_GET_ALL_ACCOUNTS)
	sendFax := path.Join(utilities.API_PATHS, API_UI_SEND_FAX)

	router.GET(authtenticationPath, routeAuthentication)
	router.POST(saveSettings, routeSaveSettings)
	router.GET(loadSettings, routeLoadSettings)
	router.GET(loadAccountInfo, routeAccountInfo)
	router.GET(logout, routeLogout)
	router.GET(getLastFaxes, routeLastFaxes)
	router.GET(loadAllAccounts, routeLoadAllAccounts)
	router.POST(sendFax, routeSendFax)
}

// routeSendFax handles the API route for sending a fax.
// It follows these steps:
// 1. Load user data from the settings file.
// 2. Authenticate the user and get an authentication token.
// 3. Parse and handle the multipart form data.
// 4. Extract relevant form values and files from the request.
// 5. Unmarshal JSON data into respective structures.
// 6. Read file contents.
// 7. Send the fax using user data, authentication token, and other relevant data.
//
// Parameters:
//   - c: Gin context for the HTTP request.
func routeSendFax(c *gin.Context) {
	err := c.Request.ParseMultipartForm(utilities.TWO_GB_SIZE)
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
		return
	}

	contactJSON := c.Request.FormValue("contact")
	documentJSON := c.Request.FormValue("document")
	transmissionJSON := c.Request.FormValue("transmission")
	fileModelJSON := c.Request.FormValue("fileModel")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file from request"})
		return
	}

	var contact Contact
	var document DocumentRecord
	var transmission Transmission
	var fileModel SendFileInfo

	if !unmarshalJSON(contactJSON, &contact, "failed to unmarshal contact data", c, http.StatusBadRequest) {
		return
	}

	if !unmarshalJSON(documentJSON, &document, "failed to unmarshal document data", c, http.StatusBadRequest) {
		return
	}

	if !unmarshalJSON(transmissionJSON, &transmission, "failed to unmarshal transmission data", c, http.StatusBadRequest) {
		return
	}

	if !unmarshalJSON(fileModelJSON, &fileModel, "failed to unmarshal fileModel data", c, http.StatusBadRequest) {
		return
	}

	fileContents, err := io.ReadAll(file)
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file contents"})
		return
	}

	directCall := NewApiServerDirectCalls()
	err = directCall.SendFax(contact, document, transmission, fileContents, fileModel)
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "fax sent successfully"})
}

// routeLoadAllAccounts handles the API route for retrieving account information.
// It follows these steps:
// 1. Load user data from the settings file.
// 2. Authenticate the user and get an authentication token.
// 3. Retrieve account information.
//
// Parameters:
//   - c: Gin context for the HTTP request.
func routeLoadAllAccounts(c *gin.Context) {
	directCall := NewApiServerDirectCalls()
	accountResponses, err := directCall.GetAllAccounts()
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accountResponses)
}

// routeLastFaxes handles the API route for fetching fax transmissions.
// It follows these steps:
// 1. Load user data from the settings file.
// 2. Authenticate the user and get an authentication token.
// 3. Retrieve fax transmissions.
// 4. Filter fax responses based on print status.
// 5. Convert and return filtered fax responses as fax data.
//
// Parameters:
//   - c: Gin context for the HTTP request.
func routeLastFaxes(c *gin.Context) {
	directCall := NewApiServerDirectCalls()
	faxeList, err := directCall.GetLastFaxes(0)
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, faxeList)
}

// routeAuthentication handles the API route for user authentication.
// It follows these steps:
// 1. Load user data from the settings file.
// 2. Authenticate the user and get an authentication token.
//
// Parameters:
//   - c: Gin context for the HTTP request.
func routeAuthentication(c *gin.Context) {
	userData, err := loadUserDataFromFile()
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "the settings file not found"})
		return
	}

	authResponse, err := AuthenticateICT(*userData)
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in get  key from api the host"})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// routeSaveSettings handles the API route for saving user settings.
// It follows these steps:
// 1. Decode user data from the request body.
// 2. Marshal user data to JSON.
// 3. Encrypt JSON data.
// 4. Save encrypted data to the settings file.
//
// Parameters:
//   - c: Gin context for the HTTP request.
func routeSaveSettings(c *gin.Context) {
	userdata := &UserData{}
	json.NewDecoder(c.Request.Body).Decode(userdata)

	directCall := NewApiServerDirectCalls()
	err := directCall.SaveSettings(*userdata)
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "ok")
}

// routeLoadSettings handles the API route for loading user settings.
// It follows these steps:
// 1. Load user data from the settings file.
// 2. Return loaded user data.
//
// Parameters:
//   - c: Gin context for the HTTP request.
func routeLoadSettings(c *gin.Context) {
	directCall := NewApiServerDirectCalls()
	userData, err := directCall.LoadSettings()
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userData)
}

// routeAccountInfo handles the API route for fetching account information.
// It follows these steps:
// 1. Load user data from the settings file.
// 2. Authenticate the user and get an authentication token.
// 3. Convert authentication response to account information.
//
// Parameters:
//   - c: Gin context for the HTTP request.
func routeAccountInfo(c *gin.Context) {
	directCall := NewApiServerDirectCalls()
	accountInfo, err := directCall.GetAccountInfo()
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accountInfo)
}

// routeLogout handles the API route for user logout.
// It follows these steps:
// 1. Get the path of the settings file.
// 2. Check if the settings file exists.
// 3. Remove the settings file.
//
// Parameters:
//   - c: Gin context for the HTTP request.
func routeLogout(c *gin.Context) {
	directCall := NewApiServerDirectCalls()
	err := directCall.Logout()
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// unmarshalJSON is a utility function to unmarshal JSON data into a target structure.
// It returns true if unmarshaling is successful, false otherwise.
//
// Parameters:
//  - data: JSON data to be unmarshaled.
//  - target: Pointer to the target structure for unmarshaling.
//  - errorMsg: Error message to be displayed in case of failure.
//  - c: Gin context for the HTTP request.
//  - statusCode: HTTP status code to be returned in case of failure.
//
// Returns:
//  - true if unmarshaling is successful, false otherwise.

func unmarshalJSON(data string, target interface{}, errorMsg string, c *gin.Context, statusCode int) bool {
	err := json.Unmarshal([]byte(data), target)
	if err != nil {
		logger.Inst().Error(err.Error())
		c.JSON(statusCode, gin.H{"error": errorMsg})
		return false
	}
	return true
}

// loadUserDataFromFile retrieves user data from the system settings file.
// It follows these steps:
// 1. Get the path of the system settings file.
// 2. Check if the settings file exists.
// 3. If the file exists, load and return the user data using getSettingsFileContent.
// 4. If the file does not exist, return an error indicating that the settings file is not found.
//
// Returns:
//   - A pointer to UserData if successful, an error otherwise.
func loadUserDataFromFile() (*UserData, error) {
	settingsFilePath, _ := utilities.GetSystemSettingsPath()
	if !utilities.CheckIfFileExists(settingsFilePath) {
		return nil, errors.New("the settings file not found")
	}
	return getSettingsFileContent()
}
