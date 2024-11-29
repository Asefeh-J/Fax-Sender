package api

import (
	"bytes"
	"encoding/json"
	"faxsender/src/utilities"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Constants defining various API paths and base URL for the ICT API.
// These paths are used to construct complete URLs for specific API endpoints.
const (
	ICT_AUTHENTICATION_API_PATH        = "api/authenticate"
	ICT_TRANSMISSION_API_PATH          = "api/transmissions"
	ICT_ACCOUNTS_API_PATH              = "api/accounts"
	ICT_CONTACTS_API_PATH              = "api/contacts"
	ICT_Document_API_PATH              = "api/documents"
	ICT_PROGRAMS_API_PATH              = "api/programs/sendfax"
	ICT_DOCUMENS_WITH_ID_API_PATH      = "api/documents/%d/media"
	ICT_TRANMISSTIONS_WITH_ID_API_PATH = "api/transmissions/%d/send"
)

// buildICTRequestURL constructs the complete URL for an ICT API endpoint.
// Steps:
// 1. Use the provided user data and URI path to construct the URL.
//
// Parameters:
//   - userData: User data containing the hostname.
//   - uriPath: URI path for the API endpoint.
//
// Returns:
//   - string: The constructed URL.
func buildICTReqeustURL(userData *UserData, uriPath string) string {
	parsedUrl, _ := url.ParseRequestURI(userData.Hostname)

	return fmt.Sprintf("%s://%s/%s", parsedUrl.Scheme, parsedUrl.Host, uriPath)
}

// AuthenticateICT performs user authentication with the ICT API.
// Steps:
// 1. Marshal user data into JSON.
// 2. Build the authentication URL.
// 3. Make an HTTP POST request to the authentication API.
// 4. Check if the authentication was successful (status code 200).
// 5. Decode the response body into an AuthResponse struct.
//
// Parameters:
//   - userData: User data containing authentication information.
//
// Returns:
//   - *AuthResponse: Authentication response containing user details.
//   - error: An error if authentication fails or any other error occurs.
func AuthenticateICT(userData UserData) (*AuthResponse, error) {

	bodyBytes, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	authURL := buildICTReqeustURL(&userData, ICT_AUTHENTICATION_API_PATH)
	resp, err := http.Post(authURL, utilities.JSON_CONTENT_TYPE, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Authentication failed with status code: %d", resp.StatusCode)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	return &authResponse, nil
}

// TransmissionsICT retrieves fax transmissions from the ICT API.
// Steps:
// 1. Build the URL for the transmissions API.
// 2. Create an authenticated HTTP GET request.
// 3. Make the HTTP request and check for success (status code 200).
// 4. Decode the response body into a slice of FaxResponse structs.
// 5. Parse the DateTime field in each response into a time.Time field.
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//
// Returns:
//   - []FaxResponse: A slice of FaxResponse structs representing fax transmissions.
//   - error: An error if fetching transmissions fails or any other error occurs.
func TransmissionsICT(userData UserData, authToken string) ([]FaxResponse, error) {

	faxURL := buildICTReqeustURL(&userData, ICT_TRANSMISSION_API_PATH)

	req, err := http.NewRequest("GET", faxURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{
		Timeout: time.Second * 10, // Set a reasonable timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Fetching faxes failed with status code: %d", resp.StatusCode)
	}

	var faxResponse []FaxResponse
	if err := json.NewDecoder(resp.Body).Decode(&faxResponse); err != nil {
		return nil, err
	}

	for i, response := range faxResponse {
		lastRunTimestamp, err := strconv.ParseInt(response.DateTime, 10, 64)
		if err == nil {
			faxResponse[i].DateTimeParsed = time.Unix(lastRunTimestamp, 0)
		}
	}

	return faxResponse, nil
}

// AccountsICT retrieves account information from the ICT API.
// Steps:
// 1. Build the URL for the accounts API.
// 2. Create an authenticated HTTP GET request.
// 3. Make the HTTP request and check for success (status code 200).
// 4. Decode the response body into a slice of AccountResponse structs.
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//
// Returns:
//   - []AccountResponse: A slice of AccountResponse structs representing account information.
//   - error: An error if fetching accounts fails or any other error occurs.
func AccountsICT(userData UserData, authToken string) ([]AccountResponse, error) {

	accountURL := buildICTReqeustURL(&userData, ICT_ACCOUNTS_API_PATH)

	req, err := http.NewRequest("GET", accountURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{
		Timeout: time.Second * 10, // Set a reasonable timeout
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Fetching accounts failed with status code: %d", resp.StatusCode)
	}

	var accountResponse []AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountResponse); err != nil {
		return nil, err
	}

	return accountResponse, nil
}

// SendFaxICT sends a fax using various ICT API endpoints.
// Steps:
// 1. Convert the AccountID in the Transmission struct to an integer.
// 2. Create a Contact, Document Record, and Program sequentially.
// 3. Upload the document file.
// 4. Create a Transmission with the provided data.
// 5. Send the created Transmission.
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//   - contact: Contact information for the fax transmission.
//   - document: Document information for the fax transmission.
//   - transmission: Transmission information, including AccountID, for the fax.
//   - fileContents: Contents of the document file to be uploaded.
//   - fileModel: Information about the document file content type.
//
// Returns:
//   - error: An error if any step of the fax transmission process fails.

func SendFaxICT(userData UserData, authToken string, contact Contact, document DocumentRecord,
	transmission Transmission, fileContents []byte, fileModel SendFileInfo) error {

	accountID, _ := strconv.Atoi(transmission.AccountID)
	contentType := fileModel.ContentType
	// Step 1: Create Contact
	contactID, err := CreateContact(userData, authToken, contact)
	if err != nil {
		return fmt.Errorf("Failed to create contact: %v", err)
	}

	// Step 2: Create Document Record
	documentID, err := CreateDocumentRecord(userData, authToken, document)
	if err != nil {
		return fmt.Errorf("Failed to create document record: %v", err)
	}

	// Step 3: Upload Document File
	err = UploadDocumentFile(userData, authToken, documentID, fileContents, contentType)
	if err != nil {
		return fmt.Errorf("Failed to upload document file: %v", err)
	}

	// Step 4: Create Program
	programID, err := CreateProgram(userData, authToken, documentID)
	if err != nil {
		return fmt.Errorf("Failed to create program: %v", err)
	}

	// Step 5: Create Transmission
	transmissionID, err := CreateTransmission(userData, authToken, transmission, contactID, accountID, programID)
	if err != nil {
		return fmt.Errorf("Failed to create transmission: %v", err)
	}

	// Step 6: Send Transmission
	err = SendTransmission(userData, authToken, transmissionID)
	if err != nil {
		return fmt.Errorf("Failed to send transmission: %v", err)
	}

	return nil
}

// CreateContact creates a contact using the ICT API.
// Steps:
// 1. Build the URL for creating a contact.
// 2. Marshal the contact data into JSON.
// 3. Make an authenticated POST request and check for success (status code 200).
// 4. Read and convert the response body (contact ID) into an integer.
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//   - contact: Contact information to be created.
//
// Returns:
//   - int: The ID of the created contact.
//   - error: An error if the creation process fails.
func CreateContact(userData UserData, authToken string, contact Contact) (int, error) {
	url := buildICTReqeustURL(&userData, ICT_CONTACTS_API_PATH)

	bodyBytes, err := json.Marshal(contact)
	if err != nil {
		return 0, fmt.Errorf("error marshaling contact data: %v", err)
	}

	resp, err := makeAuthenticatedPostRequest(url, authToken, bodyBytes)
	if err != nil {
		return 0, fmt.Errorf("error making authenticated POST request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("CreateContact API call failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	contactID, _ := strconv.Atoi(string(body))

	return contactID, nil
}

// CreateDocumentRecord creates a document record using the ICT API.
// Steps:
// 1. Build the URL for creating a document record.
// 2. Marshal the document data into JSON.
// 3. Make an authenticated POST request and check for success (status code 200).
// 4. Read and convert the response body (document ID) into an integer.
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//   - document: Document information to be created.
//
// Returns:
//   - int: The ID of the created document record.
//   - error: An error if the creation process fails.
func CreateDocumentRecord(userData UserData, authToken string, document DocumentRecord) (int, error) {

	url := buildICTReqeustURL(&userData, ICT_Document_API_PATH)

	bodyBytes, err := json.Marshal(document)

	if err != nil {
		return 0, err
	}

	resp, err := makeAuthenticatedPostRequest(url, authToken, bodyBytes)

	if err != nil {
		return 0, fmt.Errorf("error making authenticated POST request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("CreateDocumentRecord API call failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	documentID, _ := strconv.Atoi(string(body))

	return documentID, nil
}

// UploadDocumentFile uploads a document file to the ICT API.
// Steps:
// 1. Build the URL for uploading a document file.
// 2. Create a multipart form and add the file field.
// 3. Make an authenticated POST request with the multipart form.
// 4. Check for success (status code 200).
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//   - documentID: The ID of the document to which the file will be attached.
//   - fileContents: Contents of the document file to be uploaded.
//   - contentType: Content type of the document file.
//
// Returns:
//   - error: An error if the upload process fails.
func UploadDocumentFile(userData UserData, authToken string, documentID int, fileContents []byte, contentType string) error {
	documentUrl := fmt.Sprintf(ICT_DOCUMENS_WITH_ID_API_PATH, documentID)
	url := buildICTReqeustURL(&userData, documentUrl)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileField, _ := writer.CreateFormFile("file", "filename.txt")
	fileField.Write(fileContents)

	writer.Close()

	req, err := makeAuthenticatedPostRequestWithBody(url, authToken, contentType, body)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		return fmt.Errorf("Uploading document file failed with status code: %d", req.StatusCode)
	}

	return nil
}

// CreateProgram creates a program using the ICT API.
// Steps:
// 1. Build the URL for creating a program.
// 2. Marshal the program data into JSON.
// 3. Make an authenticated POST request and check for success (status code 200).
// 4. Read and convert the response body (program ID) into an integer.
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//   - documentID: The ID of the document associated with the program.
//
// Returns:
//   - int: The ID of the created program.
//   - error: An error if the creation process fails.
func CreateProgram(userData UserData, authToken string, documentID int) (int, error) {

	url := buildICTReqeustURL(&userData, ICT_PROGRAMS_API_PATH)

	bodyBytes, err := json.Marshal(map[string]int{"document_id": documentID})

	if err != nil {
		return 0, err
	}

	resp, err := makeAuthenticatedPostRequest(url, authToken, bodyBytes)

	if err != nil {
		return 0, fmt.Errorf("error making authenticated POST request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response body: %s\n", errorBody)
		return 0, fmt.Errorf("CreateProgram API call failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	programID, err := strconv.Atoi(strings.TrimSpace(string(body)))

	if err != nil {
		return 0, fmt.Errorf("error converting program ID: %v", err)
	}

	return programID, nil
}

// CreateTransmission creates a transmission using the ICT API.
// Steps:
// 1. Build the URL for creating a transmission.
// 2. Convert the contact and program IDs to strings.
// 3. Marshal the transmission data into JSON.
// 4. Make an authenticated POST request and check for success (status code 200).
// 5. Read and convert the response body (transmission ID) into an integer.
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//   - transmission: Transmission information to be created.
//   - contactID: The ID of the associated contact.
//   - accountID: The ID of the associated account.
//   - programID: The ID of the associated program.
//
// Returns:
//   - int: The ID of the created transmission.
//   - error: An error if the creation process fails.
func CreateTransmission(userData UserData, authToken string, transmission Transmission, contactID, _, programID int) (int, error) {

	url := buildICTReqeustURL(&userData, ICT_TRANSMISSION_API_PATH)

	transmission.ContactID = strconv.Itoa(contactID)
	transmission.ProgramID = strconv.Itoa(programID)

	convertedTransmission := convertTransmission(transmission)

	bodyBytes, err := json.Marshal(convertedTransmission)
	if err != nil {
		return 0, err
	}

	resp, err := makeAuthenticatedPostRequest(url, authToken, bodyBytes)

	if err != nil {
		return 0, fmt.Errorf("error making authenticated POST request: %v", err)
	}

	defer resp.Body.Close()

	fmt.Println("Transmission Request Body:", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("CreateTransmission API call failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Println("Transmission Response body is:", string(body))

	transmissionID, _ := strconv.Atoi(string(body))

	return transmissionID, nil
}

// SendTransmission sends a transmission using the ICT API.
// Steps:
// 1. Build the URL for sending a transmission.
// 2. Make an authenticated POST request and check for success (status code 200).
//
// Parameters:
//   - userData: User data containing ICT API access information.
//   - authToken: Authentication token for making authenticated requests.
//   - transmissionID: The ID of the transmission to be sent.
//
// Returns:
//   - error: An error if the sending process fails.
func SendTransmission(userData UserData, authToken string, transmissionID int) error {
	transmissionUrl := fmt.Sprintf(ICT_TRANMISSTIONS_WITH_ID_API_PATH, transmissionID)
	url := buildICTReqeustURL(&userData, transmissionUrl)

	resp, err := makeAuthenticatedPostRequest(url, authToken, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Sending transmission failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// makeAuthenticatedPostRequest makes an authenticated HTTP POST request to the provided URL.
// Steps:
// 1. Create a new HTTP POST request.
// 2. Set the necessary headers (content type and authorization).
// 3. Make the request using an HTTP client with a timeout.
// 4. Return the HTTP response.
//
// Parameters:
//   - url: The URL for the POST request.
//   - authToken: Authentication token for making authenticated requests.
//   - body: The request body.
//
// Returns:
//   - *http.Response: The HTTP response.
//   - error: An error if the request fails.
func makeAuthenticatedPostRequest(url, authToken string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return client.Do(req)
}

// makeAuthenticatedPostRequestWithBody makes an authenticated HTTP PUT request with a body.
// Steps:
// 1. Create a new HTTP PUT request.
// 2. Set the necessary headers (content type and authorization).
// 3. Make the request using an HTTP client with a timeout.
// 4. Return the HTTP response.
//
// Parameters:
//   - url: The URL for the PUT request.
//   - authToken: Authentication token for making authenticated requests.
//   - contentType: Content type of the request body.
//   - body: The request body.
//
// Returns:
//   - *http.Response: The HTTP response.
//   - error: An error if the request fails.
func makeAuthenticatedPostRequestWithBody(url, authToken, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return client.Do(req)
}
