package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"io"
)

// AESEncryptGCM encrypts the given data using AES-GCM mode.
//
// This function provides authenticated encryption using AES-GCM mode.
// The nonce is randomly generated and prepended to the output.
//
// Format: Nonce (12 bytes) || Ciphertext || Authentication Tag (16 bytes)
//
// Parameters:
//
//	data []byte - The plaintext data to be encrypted.
//	key  []byte - The AES key (must be 16, 24, or 32 bytes for AES-128/192/256).
//	aad  []byte - Additional authenticated data (AAD) to be included in the authentication tag.
//	              Can be nil if no AAD is needed.
//
// Returns:
//
//	[]byte - The encrypted output: Nonce (12 bytes) + Ciphertext + Auth Tag (16 bytes).
//	error  - An error if encryption fails or the key length is invalid.
//
// Technical Details:
//   - Nonce: 12 bytes (96 bits) - randomly generated for each encryption
//   - Authentication Tag: 16 bytes (128 bits) - automatically appended by GCM
//   - Total overhead: 28 bytes (12-byte nonce + 16-byte tag)
func AESEncryptGCM(data, key []byte, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
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
	ciphertext := gcm.Seal(nonce, nonce, data, aad)
	return ciphertext, nil
}

// AESDecryptGCM decrypts the given ciphertext using AES-GCM mode.
//
// This function performs authenticated decryption by extracting the nonce,
// verifying the authentication tag, and decrypting the ciphertext.
//
// Expected Format: Nonce (12 bytes) || Ciphertext || Authentication Tag (16 bytes)
//
// Parameters:
//
//	crypted []byte - The encrypted data with nonce prepended and auth tag appended.
//	                 Minimum length: 28 bytes (12-byte nonce + 16-byte tag).
//	key     []byte - The AES key (must be 16, 24, or 32 bytes for AES-128/192/256).
//	aad     []byte - Additional authenticated data (AAD) used during encryption.
//	                 Must match exactly what was used during encryption.
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails, authentication fails, or the key length is invalid.
//
// Technical Details:
//   - Nonce: First 12 bytes of crypted data
//   - Authentication Tag: Last 16 bytes (verified automatically by GCM)
//   - Minimum input size: 28 bytes
func AESDecryptGCM(crypted, key []byte, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(crypted) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := crypted[:nonceSize], crypted[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, aad)
}

// AESEncryptGCMPEM encrypts the given data using AES-GCM mode with a PEM-encoded key.
//
// This function decodes a PEM-formatted AES key and uses it for GCM encryption.
// The PEM block should contain the raw AES key bytes.
//
// Output Format: Nonce (12 bytes) || Ciphertext || Authentication Tag (16 bytes)
//
// Parameters:
//
//	data []byte - The plaintext data to be encrypted.
//	key  []byte - The PEM-encoded AES key (containing 16/24/32 raw bytes).
//	aad  []byte - Additional authenticated data (AAD). Can be nil.
//
// Returns:
//
//	[]byte - The encrypted output: Nonce (12 bytes) + Ciphertext + Auth Tag (16 bytes).
//	error  - An error if encryption fails, PEM decoding fails, or the key is invalid.
//
// Technical Details:
//   - Total overhead: 28 bytes (12-byte nonce + 16-byte tag)
//   - PEM block type can be any valid type containing raw AES key bytes
func AESEncryptGCMPEM(data, key []byte, aad []byte) ([]byte, error) {
	aeskey, rest := pem.Decode(key)
	if len(rest) > 0 || aeskey == nil {
		return nil, errors.New("invalid aes key")
	}
	return AESEncryptGCM(data, aeskey.Bytes, aad)
}

// AESDecryptGCMPEM decrypts the given ciphertext using AES-GCM mode with a PEM-encoded key.
//
// This function decodes a PEM-formatted AES key and uses it for GCM decryption.
// The PEM block should contain the raw AES key bytes.
//
// Expected Format: Nonce (12 bytes) || Ciphertext || Authentication Tag (16 bytes)
//
// Parameters:
//
//	data []byte - The encrypted data with nonce and auth tag (minimum 28 bytes).
//	key  []byte - The PEM-encoded AES key (containing 16/24/32 raw bytes).
//	aad  []byte - Additional authenticated data (AAD) used during encryption.
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails, PEM decoding fails, authentication fails, or the key is invalid.
//
// Technical Details:
//   - Minimum input size: 28 bytes
//   - Authentication is verified automatically during decryption
func AESDecryptGCMPEM(data, key []byte, aad []byte) ([]byte, error) {
	aeskey, rest := pem.Decode(key)
	if len(rest) > 0 || aeskey == nil {
		return nil, errors.New("invalid aes key")
	}
	return AESDecryptGCM(data, aeskey.Bytes, aad)
}
