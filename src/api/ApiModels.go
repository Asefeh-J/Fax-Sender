package api

import (
	"encoding/json"
	"faxsender/src/utilities"
	"os"
	"strconv"
	"time"
)

// Constants for API UI endpoints.
const (
	EMPTY_FIELD string = "N/A"

	API_UI_AUTHENTICATION    = "authentication"
	API_UI_SAVE_SETTINGS     = "save_settings"
	API_UI_LOAD_SETTINGS     = "load_settings"
	API_UI_LOAD_ACCOUNT_INFO = "load_account_info"
	API_UI_LOGOUT            = "logout"
	API_UI_GET_LAST_FAXES    = "load_faxes"
	API_UI_GET_ALL_ACCOUNTS  = "load_accounts"
	API_UI_SEND_FAX          = "send_fax"
)

// AccountInfo represents user account information shown on the second tab.
type AccountInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
	Company   string `json:"company"`
}

// IApiUICalls represents the interface as dependency injection for the api calls
type IApiUICalls interface {
	GetAccountInfo() (*AccountInfo, error)
	SaveSettings(userData UserData) error
	LoadSettings() (*UserData, error)
	Logout() error
	GetLastFaxes(count int) ([]FaxData, error)
	GetAllAccounts() ([]AccountResponse, error)
	SendFax(contact Contact, document DocumentRecord, transmission Transmission, file []byte, fileModel SendFileInfo) error
}

// UserData represents user credentials to log in.
type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"host"`
}

// Contact represents contact information for fax destination.
type Contact struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Custom1     string `json:"custom1"`
	Custom2     string `json:"custom2"`
	Custom3     string `json:"custom3"`
	Description string `json:"description"`
}

// DocumentRecord represents Title and Description of the fax in the process of sending fax.
type DocumentRecord struct {
	Title       string `json:"name"`
	Description string `json:"description"`
}

// Program represents a program entity used in the process of sending fax.
type Program struct {
	DocumentID string `json:"document_id"`
	File       []byte `json:"file"`
}

// Transmission represents transmission data in the process of sending fax.
type Transmission struct {
	Title       string `json:"title"`
	ContactID   string `json:"contact_id"`
	AccountID   string `json:"account_id"`
	ProgramID   string `json:"program_id"`
	IsPrint     string `json:"is_print"`
	IsCoverPage string `json:"is_coverpage"`
	TryAllowed  string `json:"try_allowed"`
}

// ConvertedTransmission represents the converted transmission data.
type ConvertedTransmission struct {
	Title       string `json:"title"`
	ContactID   int    `json:"contact_id"`
	AccountID   int    `json:"account_id"`
	ProgramID   int    `json:"program_id"`
	IsPrint     int    `json:"is_print"`
	IsCoverPage int    `json:"is_coverpage"`
	TryAllowed  int    `json:"try_allowed"`
}

// AuthResponse represents the authentication response from api call.
type AuthResponse struct {
	Token     string `json:"token"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
	Company   string `json:"company"`
}

// AccountResponse represents the response containing account details from api call.
type AccountResponse struct {
	AccountID string `json:"account_id"`
	TenantID  string `json:"tenant_id"`
	Type      string `json:"type"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedBy string `json:"created_by"`
	Company   string `json:"company"`
}

// FaxData represents fax-related data shown on the third tab.
type FaxData struct {
	DateTime       time.Time `json:"last_run"`
	Title          string    `json:"title"`
	DestinationFax string    `json:"contact_phone"`
	CallerID       string    `json:"account_phone"`
	Status         string    `json:"status"`
}

// FaxResponse represents the response containing fax details from api call.
type FaxResponse struct {
	DateTime       string `json:"last_run"`
	Title          string `json:"title"`
	DestinationFax string `json:"contact_phone"`
	CallerID       string `json:"account_phone"`
	Status         string `json:"status"`
	Is_Print       string `json:"is_print"`
	DateTimeParsed time.Time
}

// SendFileInfo represents information about the file Content-Type.
type SendFileInfo struct {
	ContentType string `json:"content_type"`
}

// getSettingsFileContent reads and decrypts user settings from the file.
//
// Steps:
// 1. Retrieve the path of the system settings file using utilities.GetSystemSettingsPath.
// 2. Read the contents of the settings file using os.ReadFile.
// 3. Decrypt the read content using utilities.DecryptData.
// 4. Unmarshal the decrypted content into a UserData struct using json.Unmarshal.
//
// Parameters:
//   - None
//
// Returns:
//   - *UserData: A pointer to the UserData struct containing the decrypted user settings.
//   - error: An error, if any, encountered during the process of reading, decrypting,
//     or unmarshalling the user settings file.
func getSettingsFileContent() (*UserData, error) {
	settingsFilePath, _ := utilities.GetSystemSettingsPath()

	settingsFileContent, err := os.ReadFile(settingsFilePath)
	if err != nil {
		return nil, err
	}

	decryptedSettings, err := utilities.DecryptData(settingsFileContent)
	if err != nil {
		return nil, err
	}

	var userData *UserData
	err = json.Unmarshal(decryptedSettings, &userData)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

// GetEmptyAccountInfo returns a pre-defined empty account information.
//
// Steps:
// 1. Create an instance of the AccountInfo struct with pre-defined empty values.
//
// Parameters:
//
//	None
//
// Returns:
//   - *AccountInfo: A pointer to the AccountInfo struct with pre-defined empty values.
func GetEmptyAccountInfo() *AccountInfo {
	return &AccountInfo{
		FirstName: EMPTY_FIELD,
		LastName:  EMPTY_FIELD,
		Email:     "empty@empty.com",
		Phone:     "+(0) 000 000-000",
		Country:   EMPTY_FIELD,
		Company:   EMPTY_FIELD,
	}
}

// ConvertAuthResponseToAccountInfo converts AuthResponse to AccountInfo.
//
// Steps:
// 1. Create an instance of the AccountInfo struct using values from the provided AuthResponse.
//
// Parameters:
//   - res: AuthResponse containing authentication response data.
//
// Returns:
//   - *AccountInfo: A pointer to the AccountInfo struct with values converted from the AuthResponse.
func ConvertAuthResponseToAccountInfo(res AuthResponse) *AccountInfo {

	return &AccountInfo{
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Phone:     res.Phone,
		Country:   EMPTY_FIELD,
		Company:   res.Company,
	}
}

// ConvertFilteredFaxResponsesToFaxData converts filtered fax responses to FaxData.
//
// Steps:
// 1. Create a slice of FaxData with the same length as the provided filtered fax responses.
// 2. Iterate over each FaxResponse in the provided list, converting each response to a FaxData instance.
// 3. Parse the DateTime field of each FaxResponse to a Unix timestamp and convert it to a time.Time object.
// 4. Populate the FaxData fields with values from the corresponding FaxResponse.
// 5. Append the converted FaxData to the result slice.
//
// Parameters:
//   - res: A slice of FaxResponse containing filtered fax response data.
//
// Returns:
//   - []FaxData: A slice of FaxData instances converted from the filtered fax responses.
func ConvertFilteredFaxResponsesToFaxData(res []FaxResponse) []FaxData {
	faxDataList := make([]FaxData, len(res))

	for i, response := range res {
		lastRunTimestamp, err := strconv.ParseInt(response.DateTime, 10, 64)
		if err == nil {
			lastRunTime := time.Unix(lastRunTimestamp, 0)
			faxDataList[i] = FaxData{
				DateTime:       lastRunTime,
				Title:          response.Title,
				DestinationFax: response.DestinationFax,
				CallerID:       response.CallerID,
				Status:         response.Status,
			}
		}
	}

	return faxDataList
}

// convertTransmission converts Transmission to ConvertedTransmission (string fields to int fields).
//
// Steps:
// 1. Convert string fields of the provided Transmission to their corresponding int values.
// 2. Create an instance of ConvertedTransmission using the converted int values.
//
// Parameters:
//   - transmission: The original Transmission instance with string fields.
//
// Returns:
//   - ConvertedTransmission: An instance of ConvertedTransmission with int fields.
func convertTransmission(transmission Transmission) ConvertedTransmission {
	contactID, _ := strconv.Atoi(transmission.ContactID)
	programID, _ := strconv.Atoi(transmission.ProgramID)
	accountID, _ := strconv.Atoi(transmission.AccountID)
	isPrint, _ := strconv.Atoi(transmission.IsPrint)
	isCoverPage, _ := strconv.Atoi(transmission.IsCoverPage)
	tryAllowed, _ := strconv.Atoi(transmission.TryAllowed)

	converted := ConvertedTransmission{
		Title:       transmission.Title,
		ContactID:   contactID,
		AccountID:   accountID,
		ProgramID:   programID,
		IsPrint:     isPrint,
		IsCoverPage: isCoverPage,
		TryAllowed:  tryAllowed,
	}

	return converted
}
