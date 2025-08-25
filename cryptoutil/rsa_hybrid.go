package cryptoutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
)

// RSAHybridEncrypt_ encrypts data using RSA-OAEP + AES-GCM hybrid encryption.
//
// This function provides secure encryption for data of any size by combining
// RSA-OAEP (for key encryption) with AES-GCM (for data encryption).
//
// Output Format: RSA Key Length (4 bytes) || Encrypted AES Key || Encrypted Data
//
// Parameters:
//
//	data   []byte - The plaintext data to be encrypted (any size).
//	pubKey *rsa.PublicKey - The RSA public key for encrypting the AES key.
//	aad    []byte - Additional authenticated data for AES-GCM. Can be nil.
//
// Returns:
//
//	[]byte - The hybrid encrypted data.
//	error  - An error if encryption fails.
//
// Technical Details:
//   - AES Key: 32 bytes (AES-256) - randomly generated
//   - RSA Encryption: OAEP with SHA-256
//   - AES Encryption: GCM mode with 12-byte nonce and 16-byte tag
//   - Format: [4-byte key length][RSA-encrypted AES key][AES-GCM encrypted data]
func RSAHybridEncrypt_(data []byte, pubKey *rsa.PublicKey, aad []byte) ([]byte, error) {
	// 1. Generate random AES-256 key
	aesKey := GenerateAESKey() // 32 bytes

	// 2. Encrypt data with AES-GCM
	encryptedData, err := AESEncryptGCM(data, aesKey, aad)
	if err != nil {
		return nil, err
	}

	// 3. Encrypt AES key with RSA-OAEP
	// Note: Single RSA operation since AES key (32 bytes) fits within RSA limits
	keySize := pubKey.Size()
	maxKeySize := keySize - 2*32 - 2 // OAEP overhead: 2*SHA256_SIZE + 2 = 66 bytes
	if len(aesKey) > maxKeySize {
		return nil, errors.New("AES key too large for RSA key size")
	}

	encryptedKey, err := RSAEncrypt_([]byte(aesKey), pubKey)
	if err != nil {
		return nil, err
	}

	// 4. Combine results: [key_length][encrypted_key][encrypted_data]
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, uint32(len(encryptedKey)))
	result = append(result, encryptedKey...)
	result = append(result, encryptedData...)

	return result, nil
}

// RSAHybridDecrypt_ decrypts data using RSA-OAEP + AES-GCM hybrid decryption.
//
// This function decrypts hybrid encrypted data by first extracting and
// decrypting the AES key with RSA-OAEP, then decrypting the data with AES-GCM.
//
// Expected Format: RSA Key Length (4 bytes) || Encrypted AES Key || Encrypted Data
//
// Parameters:
//
//	crypted []byte - The hybrid encrypted data.
//	priKey  *rsa.PrivateKey - The RSA private key for decrypting the AES key.
//	aad     []byte - Additional authenticated data used during encryption.
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails or authentication fails.
//
// Technical Details:
//   - Minimum input size: 32 bytes (4-byte length + RSA key + AES-GCM overhead)
//   - RSA Decryption: OAEP with SHA-256
//   - AES Decryption: GCM mode with authentication verification
func RSAHybridDecrypt_(crypted []byte, priKey *rsa.PrivateKey, aad []byte) ([]byte, error) {
	if len(crypted) < 4 {
		return nil, errors.New("ciphertext too short")
	}

	// 1. Extract encrypted key length
	keyLength := binary.BigEndian.Uint32(crypted[:4])
	if len(crypted) < int(4+keyLength) {
		return nil, errors.New("invalid ciphertext format")
	}

	// 2. Extract encrypted AES key and encrypted data
	encryptedKey := crypted[4 : 4+keyLength]
	encryptedData := crypted[4+keyLength:]

	// 3. Decrypt AES key with RSA-OAEP
	aesKey, err := RSADecrypt_(encryptedKey, priKey)
	if err != nil {
		return nil, err
	}

	// 4. Decrypt data with AES-GCM
	plaintext, err := AESDecryptGCM(encryptedData, aesKey, aad)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// RSAHybridEncrypt encrypts data using RSA-OAEP + AES-GCM hybrid encryption with byte key.
//
// This function provides secure encryption for data of any size using hybrid cryptography.
//
// Parameters:
//
//	data   []byte - The plaintext data to be encrypted (any size).
//	pubKey []byte - The RSA public key in PKCS#1 format.
//	aad    []byte - Additional authenticated data for AES-GCM. Can be nil.
//
// Returns:
//
//	[]byte - The hybrid encrypted data.
//	error  - An error if encryption fails or key parsing fails.
func RSAHybridEncrypt(data, pubKey []byte, aad []byte) ([]byte, error) {
	pub, err := x509.ParsePKCS1PublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return RSAHybridEncrypt_(data, pub, aad)
}

// RSAHybridDecrypt decrypts data using RSA-OAEP + AES-GCM hybrid decryption with byte key.
//
// This function decrypts hybrid encrypted data using RSA-OAEP + AES-GCM.
//
// Parameters:
//
//	crypted []byte - The hybrid encrypted data.
//	priKey  []byte - The RSA private key in PKCS#1 format.
//	aad     []byte - Additional authenticated data used during encryption.
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails, key parsing fails, or authentication fails.
func RSAHybridDecrypt(crypted, priKey []byte, aad []byte) ([]byte, error) {
	pri, err := x509.ParsePKCS1PrivateKey(priKey)
	if err != nil {
		return nil, err
	}
	return RSAHybridDecrypt_(crypted, pri, aad)
}

// RSAHybridEncryptPEM encrypts data using RSA-OAEP + AES-GCM hybrid encryption with PEM key.
//
// This function provides secure encryption for data of any size using PEM-formatted keys.
//
// Parameters:
//
//	data   []byte - The plaintext data to be encrypted (any size).
//	pubKey []byte - The RSA public key in PEM format.
//	aad    []byte - Additional authenticated data for AES-GCM. Can be nil.
//
// Returns:
//
//	[]byte - The hybrid encrypted data.
//	error  - An error if encryption fails, PEM decoding fails, or key parsing fails.
//
// Expected PEM Format:
//
//	-----BEGIN RSA PUBLIC KEY-----
//	[base64-encoded public key]
//	-----END RSA PUBLIC KEY-----
func RSAHybridEncryptPEM(data, pubKey []byte, aad []byte) ([]byte, error) {
	key, rest := pem.Decode(pubKey)
	if len(rest) > 0 || key == nil {
		return nil, errors.New("invalid public key")
	}
	return RSAHybridEncrypt(data, key.Bytes, aad)
}

// RSAHybridDecryptPEM decrypts data using RSA-OAEP + AES-GCM hybrid decryption with PEM key.
//
// This function decrypts hybrid encrypted data using PEM-formatted keys.
//
// Parameters:
//
//	crypted []byte - The hybrid encrypted data.
//	priKey  []byte - The RSA private key in PEM format.
//	aad     []byte - Additional authenticated data used during encryption.
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails, PEM decoding fails, key parsing fails, or authentication fails.
//
// Expected PEM Format:
//
//	-----BEGIN RSA PRIVATE KEY-----
//	[base64-encoded private key]
//	-----END RSA PRIVATE KEY-----
func RSAHybridDecryptPEM(crypted, priKey []byte, aad []byte) ([]byte, error) {
	key, rest := pem.Decode(priKey)
	if len(rest) > 0 || key == nil {
		return nil, errors.New("invalid private key")
	}
	return RSAHybridDecrypt(crypted, key.Bytes, aad)
}

// RSAHybridEncryptWithoutAAD encrypts data using RSA-OAEP + AES-GCM hybrid encryption without AAD.
//
// This is a convenience function for hybrid encryption without additional authenticated data.
//
// Parameters:
//
//	data   []byte - The plaintext data to be encrypted (any size).
//	pubKey []byte - The RSA public key in PKCS#1 format.
//
// Returns:
//
//	[]byte - The hybrid encrypted data.
//	error  - An error if encryption fails or key parsing fails.
func RSAHybridEncryptWithoutAAD(data, pubKey []byte) ([]byte, error) {
	return RSAHybridEncrypt(data, pubKey, nil)
}

// RSAHybridDecryptWithoutAAD decrypts data using RSA-OAEP + AES-GCM hybrid decryption without AAD.
//
// This is a convenience function for hybrid decryption without additional authenticated data.
//
// Parameters:
//
//	crypted []byte - The hybrid encrypted data.
//	priKey  []byte - The RSA private key in PKCS#1 format.
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails, key parsing fails, or authentication fails.
func RSAHybridDecryptWithoutAAD(crypted, priKey []byte) ([]byte, error) {
	return RSAHybridDecrypt(crypted, priKey, nil)
}
