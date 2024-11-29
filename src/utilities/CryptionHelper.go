package utilities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// getEncryptionKey returns the encryption key used for AES-GCM encryption.
func getEncryptionKey() []byte {
	key := []byte(ENCRYPTION_KEY)
	return key
}

// EncryptData encrypts binary data using AES-GCM encryption.
//
// Steps:
// 1. Create a new AES cipher block using the encryption key.
// 2. Create a new AES-GCM instance for authenticated encryption.
// 3. Generate a random nonce of the appropriate size.
// 4. Encrypt the binary data using the nonce and the AES-GCM instance.
// 5. Return the ciphertext.
//
// Parameters:
//   - binaryData: Binary data to be encrypted.
//
// Returns:
//   - []byte: The ciphertext of the encrypted data.
//   - error: An error if the encryption process fails.
func EncryptData(binaryData []byte) ([]byte, error) {
	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, binaryData, nil)
	return ciphertext, nil
}

// DecryptData decrypts encrypted data using AES-GCM decryption.
//
// Steps:
// 1. Create a new AES cipher block using the encryption key.
// 2. Create a new AES-GCM instance for authenticated decryption.
// 3. Extract the nonce and ciphertext from the encrypted data.
// 4. Decrypt the ciphertext using the nonce and the AES-GCM instance.
// 5. Return the plaintext.
//
// Parameters:
//   - encryptedData: Encrypted data to be decrypted.
//
// Returns:
//   - []byte: The plaintext of the decrypted data.
//   - error: An error if the decryption process fails.
func DecryptData(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("encrypted data is too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
