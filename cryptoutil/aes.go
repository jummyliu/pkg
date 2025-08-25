package cryptoutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"io"
)

// GenerateAESKey generates a 256-bit (32 bytes) AES key using crypto/rand.
//
// This function creates a cryptographically secure random AES-256 key
// suitable for encryption operations.
//
// Returns:
//
//	[]byte - A 32-byte AES-256 key generated using crypto/rand.
//
// Technical Details:
//   - Key length: 32 bytes (256 bits) for AES-256
//   - Uses crypto/rand.Reader for cryptographically secure randomness
func GenerateAESKey() (key []byte) {
	return GenerateAESKeyWithSize(256)
}

// GeneratePEMAESKey generates a 256-bit AES key in PEM format using crypto/rand.
//
// This function creates a cryptographically secure random AES-256 key and
// encodes it in PEM format with the block type "AES KEY".
//
// Returns:
//
//	[]byte - A PEM-encoded AES-256 key.
//
// PEM Format:
//
//	-----BEGIN AES KEY-----
//	[base64-encoded 32-byte key]
//	-----END AES KEY-----
func GeneratePEMAESKey() (key []byte) {
	aesKey := GenerateAESKey()
	key = pem.EncodeToMemory(&pem.Block{
		Type:  "AES KEY",
		Bytes: aesKey,
	})
	return key
}

// GenerateAESKeyWithSize generates an AES key of the specified bit size using crypto/rand.
//
// This function creates a cryptographically secure random AES key of the
// specified bit length. Invalid bit sizes are automatically corrected to 256 bits.
//
// Parameters:
//
//	bits int - The desired key size in bits (128, 192, or 256 for AES-128/192/256).
//	           Invalid bit sizes default to 256 bits.
//
// Returns:
//
//	[]byte - An AES key of the specified size generated using crypto/rand.
//
// Supported Bit Sizes:
//   - 128 bits: AES-128 (16 bytes)
//   - 192 bits: AES-192 (24 bytes)
//   - 256 bits: AES-256 (32 bytes) - default for invalid sizes
func GenerateAESKeyWithSize(bits int) (key []byte) {
	if bits != 128 && bits != 192 && bits != 256 {
		bits = 256
	}
	key = make([]byte, bits/8)
	_, _ = io.ReadFull(rand.Reader, key)
	return key
}

// GeneratePEMAESKeyWithSize generates an AES key of the specified bit size in PEM format.
//
// This function creates a cryptographically secure random AES key of the
// specified bit length and encodes it in PEM format with the block type "AES KEY".
//
// Parameters:
//
//	bits int - The desired key size in bits (128, 192, or 256 for AES-128/192/256).
//	           Invalid bit sizes default to 256 bits.
//
// Returns:
//
//	[]byte - A PEM-encoded AES key of the specified size.
//
// PEM Format:
//
//	-----BEGIN AES KEY-----
//	[base64-encoded key bytes]
//	-----END AES KEY-----
//
// Key Sizes:
//   - 128 bits: 16-byte key
//   - 192 bits: 24-byte key
//   - 256 bits: 32-byte key (default)
func GeneratePEMAESKeyWithSize(bits int) (key []byte) {
	aesKey := GenerateAESKeyWithSize(bits)
	key = pem.EncodeToMemory(&pem.Block{
		Type:  "AES KEY",
		Bytes: aesKey,
	})
	return key
}

// PKCS5Padding applies PKCS#5 padding to the input data.
//
// This function pads the input data to make its length a multiple of the
// specified block size. The padding bytes contain the padding length.
//
// Parameters:
//
//	data      []byte - The data to be padded.
//	blockSize int    - The block size for padding (typically 16 for AES).
//
// Returns:
//
//	[]byte - The padded data.
//
// Padding Details:
//   - Each padding byte contains the number of padding bytes added
//   - Always adds at least 1 byte of padding (1-blockSize bytes total)
//   - Example: For 16-byte blocks, if data needs 3 bytes, adds [3,3,3]
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS5UnPadding removes PKCS#5 padding from the input data.
//
// This function validates and removes PKCS#5 padding from the input data.
// The last byte indicates the number of padding bytes to remove.
//
// Parameters:
//
//	data []byte - The padded data to be unpadded.
//
// Returns:
//
//	[]byte - The unpadded data.
//	error  - An error if the padding is invalid or corrupted.
//
// Validation:
//   - Checks if padding length is valid (not exceeding data length)
//   - Empty data is returned as-is
func PKCS5UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return data, nil
	}
	unpadding := int(data[length-1])
	if length-unpadding < 0 {
		return nil, errors.New("pkcs5/unpadding failure: length - unpadding < 0")
	}
	return data[:length-unpadding], nil
}

// AESEncrypt encrypts data using AES-CBC mode with a fixed IV.
//
// ⚠️ WARNING: This function uses the first 16 bytes of the key as the IV,
// which is cryptographically insecure. Use AESEncryptCBC_HMAC instead.
//
// Parameters:
//
//	data []byte - The plaintext data to be encrypted.
//	key  []byte - The AES key (must be 16, 24, or 32 bytes).
//
// Returns:
//
//	[]byte - The encrypted ciphertext (without IV).
//	error  - An error if encryption fails or the key length is invalid.
//
// Security Issues:
//   - Uses key[:16] as IV, which is predictable and reused
//   - No authentication/integrity protection
//   - Vulnerable to chosen-ciphertext attacks
//
// Recommended Alternative: Use AESEncryptCBC_HMAC for secure encryption
func AESEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	if len(data)%blockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	return crypted, nil
}

// AESDecrypt decrypts data using AES-CBC mode with a fixed IV.
//
// ⚠️ WARNING: This function uses the first 16 bytes of the key as the IV,
// which is cryptographically insecure. Use AESDecryptCBC_HMAC instead.
//
// Parameters:
//
//	crypted []byte - The ciphertext to be decrypted.
//	key     []byte - The AES key (must be 16, 24, or 32 bytes).
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails, padding is invalid, or the key length is invalid.
//
// Security Issues:
//   - Uses key[:16] as IV, matching the insecure encryption
//   - No authentication/integrity verification
//   - Vulnerable to padding oracle attacks
//
// Recommended Alternative: Use AESDecryptCBC_HMAC for secure decryption
func AESDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(crypted)%blockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	data := make([]byte, len(crypted))
	blockMode.CryptBlocks(data, crypted)
	data, err = PKCS5UnPadding(data)
	return data, err
}

// AESEncryptPEM encrypts data using AES-CBC mode with a PEM-encoded key.
//
// ⚠️ WARNING: This function inherits the security issues from AESEncrypt.
// Use AESEncryptCBC_HMAC_PEM for secure encryption.
//
// Parameters:
//
//	data []byte - The plaintext data to be encrypted.
//	key  []byte - The PEM-encoded AES key.
//
// Returns:
//
//	[]byte - The encrypted ciphertext (without IV).
//	error  - An error if encryption fails, PEM decoding fails, or the key is invalid.
//
// PEM Format Expected:
//
//	-----BEGIN AES KEY-----
//	[base64-encoded key bytes]
//	-----END AES KEY-----
//
// Recommended Alternative: Use AESEncryptCBC_HMAC_PEM for secure encryption
func AESEncryptPEM(data, key []byte) ([]byte, error) {
	aeskey, rest := pem.Decode(key)
	if len(rest) > 0 || aeskey == nil {
		return nil, errors.New("invalid aes key")
	}
	return AESEncrypt(data, aeskey.Bytes)
}

// AESDecryptPEM decrypts data using AES-CBC mode with a PEM-encoded key.
//
// ⚠️ WARNING: This function inherits the security issues from AESDecrypt.
// Use AESDecryptCBC_HMAC_PEM for secure decryption.
//
// Parameters:
//
//	data []byte - The ciphertext to be decrypted.
//	key  []byte - The PEM-encoded AES key.
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails, PEM decoding fails, padding is invalid, or the key is invalid.
//
// PEM Format Expected:
//
//	-----BEGIN AES KEY-----
//	[base64-encoded key bytes]
//	-----END AES KEY-----
//
// Recommended Alternative: Use AESDecryptCBC_HMAC_PEM for secure decryption
func AESDecryptPEM(data, key []byte) ([]byte, error) {
	aeskey, rest := pem.Decode(key)
	if len(rest) > 0 || aeskey == nil {
		return nil, errors.New("invalid aes key")
	}
	return AESDecrypt(data, aeskey.Bytes)
}
