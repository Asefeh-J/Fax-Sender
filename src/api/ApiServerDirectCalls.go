package api

import (
	"encoding/json"
	"errors"
	"faxsender/src/utilities"
	"os"
)

// ApiUIDirectCalls represents the interface as dependency injection for the api calls without local server
type ApiServerDirectCalls struct {
	IApiUICalls
}

func NewApiServerDirectCalls() IApiUICalls {
	return &ApiServerDirectCalls{}
}

func (c *ApiServerDirectCalls) GetAccountInfo() (*AccountInfo, error) {
	userData, err := loadUserDataFromFile()
	if err != nil {
		return nil, errors.New("the settings file not found")
	}

	authResponse, err := AuthenticateICT(*userData)
	if err != nil {
		return nil, errors.New("error in get  key from api the host")
	}

	return ConvertAuthResponseToAccountInfo(*authResponse), nil
}

func (c *ApiServerDirectCalls) SaveSettings(userData UserData) error {
	jsonData, err := json.Marshal(userData)
	if err != nil {
		return errors.New("Error marshaling UserData to JSON")
	}

	encryptedData, err := utilities.EncryptData([]byte(jsonData))
	if err != nil {
		return errors.New("Error encrypting binary data")
	}

	settingsFilePath, _ := utilities.GetSystemSettingsPath()

	err = os.WriteFile(settingsFilePath, []byte(encryptedData), 0644)

	if err != nil {
		return errors.New("Error in saving data")
	}

	return nil
}

func (c *ApiServerDirectCalls) LoadSettings() (*UserData, error) {
	userData, err := loadUserDataFromFile()
	if err != nil {
		return nil, errors.New("error in load data from settings file")
	}
	return userData, nil
}

func (c *ApiServerDirectCalls) Logout() error {
	settingsFilePath, _ := utilities.GetSystemSettingsPath()
	if utilities.CheckIfFileExists(settingsFilePath) {
		err := os.Remove(settingsFilePath)
		if err != nil {
			return errors.New("error in logout!")
		}
	}
	return nil
}

func (c *ApiServerDirectCalls) GetLastFaxes(count int) ([]FaxData, error) {
	userData, err := loadUserDataFromFile()
	if err != nil {
		return nil, errors.New("error in load data from settings file")
	}

	authResponse, err := AuthenticateICT(*userData)
	if err != nil {
		return nil, errors.New("error in get  key from api the host")
	}

	authToken := authResponse.Token
	faxResponses, err := TransmissionsICT(*userData, authToken)
	if err != nil {
		return nil, errors.New("error fetching fax transmissions")
	}

	filteredFaxResponses := make([]FaxResponse, 0)
	for _, response := range faxResponses {
		if response.Is_Print == utilities.WITH_PRINT {
			filteredFaxResponses = append(filteredFaxResponses, response)
		}
	}
	startIndex := len(filteredFaxResponses) - count
	if startIndex < 0 {
		startIndex = 0
	}

	filteredFaxResponses = filteredFaxResponses[startIndex:]
	return ConvertFilteredFaxResponsesToFaxData(filteredFaxResponses), nil
}

func (c *ApiServerDirectCalls) GetAllAccounts() ([]AccountResponse, error) {
	userData, err := loadUserDataFromFile()
	if err != nil {
		return nil, errors.New("error in load data from settings file")
	}

	authResponse, err := AuthenticateICT(*userData)
	if err != nil {
		return nil, errors.New("error in get  key from api the host")
	}

	authToken := authResponse.Token
	accountResponses, err := AccountsICT(*userData, authToken)
	if err != nil {
		return nil, errors.New("error fetching fax transmissions")
	}
	return accountResponses, nil
}

func (c *ApiServerDirectCalls) SendFax(contact Contact, document DocumentRecord, transmission Transmission, fileContents []byte, fileModel SendFileInfo) error {
	userData, err := loadUserDataFromFile()
	if err != nil {
		return errors.New("error in load data from settings file")
	}

	authResponse, err := AuthenticateICT(*userData)
	if err != nil {
		return errors.New("error in getting  key from api the host")
	}
	authToken := authResponse.Token

	err = SendFaxICT(*userData, authToken, contact, document, transmission, fileContents, fileModel)
	if err != nil {
		return errors.New("error sending the fax")
	}
	return nil
}
