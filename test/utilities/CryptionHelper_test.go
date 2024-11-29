package utilities

import (
	"faxsender/src/utilities"
	"testing"
)

var (
	str string = "data to test"
)

func TestEncryptionData(t *testing.T) {
	encrypted_data, err := utilities.EncryptData([]byte(str))
	if err != nil {
		t.Error(err)
	}

	println("the file size is: ", len(encrypted_data))
}

func TestDecryptionData(t *testing.T) {
	encrypted_data, err := utilities.EncryptData([]byte(str))
	if err != nil {
		t.Error(err)
	}

	decrtyped_data, err := utilities.DecryptData(encrypted_data)
	if err != nil {
		t.Error(err)
	}

	println("the decrypted data is: ", string(decrtyped_data))
	if string(decrtyped_data) != str {
		t.Errorf("wrong decrypted data")
	}
}
