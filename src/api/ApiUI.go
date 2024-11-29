package api

import (
	"bytes"
	"encoding/json"
	"faxsender/src/utilities"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// ApiUI represents the configuration for the API server.
type ApiUI struct {
	Port int

	IApiUICalls
}

// NewApiUI creates a new ApiUI instance with the specified port.
// Steps:
// 1. Create a new ApiUI instance.
// 2. Set the port for the ApiUI instance.
//
// Parameters:
//   - port: port to connect
//
// Returns:
//   - the new instance of ApiUI
func NewApiUI(port int) IApiUICalls {
	return &ApiUI{
		Port: port,
	}
}

// buildUrl constructs the complete URL for the API endpoint.
// Steps:
// 1. Join various components to create the complete URL for the API endpoint.
//
// Parameters:
//   - endPoint: endpoint of the API
//
// Returns:
//   - the complete URL
func (a *ApiUI) buildUrl(endPoint string) string {
	return utilities.UrlJoin(
		utilities.HTTP_SCHEMA,
		utilities.LOCALHOST,
		a.Port,
		utilities.API_PATHS,
		endPoint,
	)
}

// readBody reads and parses the response body into the specified struct.
// Steps:
// 1. Ensure the response body is closed after the function exits.
// 2. Check if the API call was successful (status code 200).
// 3. Read the body of the response.
// 4. Unmarshal the body into the provided struct.
//
// Parameters:
//   - resp: HTTP response object
//   - to: destination struct for unmarshaling
//
// Returns:
//   - error if any
func (a *ApiUI) readBody(resp *http.Response, to interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &to)
	if err != nil {
		return err
	}

	return nil
}

// GetAccountInfo retrieves account information from the API.
// Steps:
// 1. Build the URL for the API endpoint.
// 2. Make an HTTP GET request to the API.
// 3. Read and parse the response body into an AccountInfo struct.
//
// Returns:
//   - AccountInfo struct
//   - error if any
func (a *ApiUI) GetAccountInfo() (*AccountInfo, error) {
	url := a.buildUrl(API_UI_LOAD_ACCOUNT_INFO)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var accountInfo *AccountInfo = &AccountInfo{}
	err = a.readBody(resp, accountInfo)
	if err != nil {
		return nil, err
	}
	return accountInfo, nil
}

// SaveSettings saves user settings via the API.
// Steps:
// 1. Build the URL for the API endpoint.
// 2. Marshal user data into JSON.
// 3. Make an HTTP POST request to the API with user data in the request body.
// 4. Check if the API call was successful (status code 200).
//
// Parameters:
//   - userData: user data to be saved
//
// Returns:
//   - error if any
func (a *ApiUI) SaveSettings(userData UserData) error {
	url := a.buildUrl(API_UI_SAVE_SETTINGS)

	data, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "content-type:"+utilities.JSON_CONTENT_TYPE, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// LoadSettings retrieves user settings from the API.
// Steps:
// 1. Build the URL for the API endpoint.
// 2. Make an HTTP GET request to the API.
// 3. Read and parse the response body into a UserData struct.
//
// Returns:
//   - UserData struct
//   - error if any
func (a *ApiUI) LoadSettings() (*UserData, error) {
	url := a.buildUrl(API_UI_LOAD_SETTINGS)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var userData *UserData = &UserData{}
	a.readBody(resp, userData)
	return userData, nil
}

// Logout performs a logout action via the API.
// Steps:
// 1. Build the URL for the API endpoint.
// 2. Make an HTTP GET request to the API.
//
// Returns:
//   - error if any
func (a *ApiUI) Logout() error {
	url := a.buildUrl(API_UI_LOGOUT)

	_, err := http.Get(url)
	return err
}

// GetLastFaxes retrieves the last N faxes from the API.
// Steps:
// 1. Build the URL for the API endpoint.
// 2. Make an HTTP GET request to the API.
// 3. Read and parse the response body into a slice of FaxData.
// 4. Retrieve the last N faxes from the slice.
//
// Parameters:
//   - count: number of faxes to retrieve
//
// Returns:
//   - slice of FaxData
//   - error if any
func (a *ApiUI) GetLastFaxes(count int) ([]FaxData, error) {
	url := a.buildUrl(API_UI_GET_LAST_FAXES)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var allFaxes []FaxData
	err = a.readBody(resp, &allFaxes)
	if err != nil {
		return nil, err
	}

	startIndex := len(allFaxes) - count
	if startIndex < 0 {
		startIndex = 0
	}

	lastFaxes := allFaxes[startIndex:]
	return lastFaxes, nil

}

// GetAllAccounts retrieves information about all accounts via the API.
// Steps:
// 1. Build the URL for the API endpoint.
// 2. Make an HTTP GET request to the API.
// 3. Read and parse the response body into a slice of AccountResponse.
//
// Returns:
//   - slice of AccountResponse
//   - error if any
func (a *ApiUI) GetAllAccounts() ([]AccountResponse, error) {
	url := a.buildUrl(API_UI_GET_ALL_ACCOUNTS)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var allAccounts []AccountResponse
	err = a.readBody(resp, &allAccounts)
	if err != nil {
		return nil, err
	}
	return allAccounts, nil
}

// SendFax sends a fax via the API.
// Steps:
// 1. Build the URL for the API endpoint.
// 2. Create a multipart form with various fields and file attachment.
// 3. Make an HTTP POST request to the API with the multipart form.
// 4. Check if the API call was successful (status code 200).
//
// Parameters:
//   - contact: Contact information for the fax.
//   - document: DocumentRecord containing details about the document to be faxed.
//   - transmission: Transmission details for the fax.
//   - file: Byte array containing the contents of the file to be faxed.
//   - fileModel: SendFileInfo providing information about the file content type.
//
// Returns:
//   - error if any
func (a *ApiUI) SendFax(contact Contact, document DocumentRecord, transmission Transmission, file []byte, fileModel SendFileInfo) error {
	url := a.buildUrl(API_UI_SEND_FAX)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := addFormField(writer, "contact", contact); err != nil {
		return err
	}

	if err := addFormField(writer, "document", document); err != nil {
		return err
	}

	if err := addFormField(writer, "transmission", transmission); err != nil {
		return err
	}

	if err := addFormField(writer, "fileModel", fileModel); err != nil {
		return err
	}

	if err := addFileField(writer, "file", "filename.txt", file); err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("APIUI call failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// addFormField adds a form field to the multipart request.
// Steps:
// 1. Create a form field in the multipart request.
// 2. Marshal the data into JSON.
// 3. Write the marshaled data to the form field.
//
// Parameters:
//   - writer: Multipart writer for creating form fields.
//   - fieldName: Name of the form field.
//   - data: Data to be added to the form field.
//
// Returns:
//   - error if any
func addFormField(writer *multipart.Writer, fieldName string, data interface{}) error {
	field, err := writer.CreateFormField(fieldName)
	if err != nil {
		return fmt.Errorf("error creating %s form field: %v", fieldName, err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling %s JSON: %v", fieldName, err)
	}

	field.Write(jsonData)
	return nil
}

// addFileField adds a file field to the multipart request.
// Steps:
// 1. Create a file field in the multipart request.
// 2. Write the file content to the file field.
//
// Parameters:
//   - writer: Multipart writer for creating file fields.
//   - fieldName: Name of the file field.
//   - fileName: Name of the file.
//   - file: Byte array containing the file content.
//
// Returns:
//   - error if any
func addFileField(writer *multipart.Writer, fieldName, fileName string, file []byte) error {
	field, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return fmt.Errorf("error creating %s form field: %v", fieldName, err)
	}

	field.Write(file)
	return nil
}
