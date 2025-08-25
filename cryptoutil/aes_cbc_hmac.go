package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/pem"
	"errors"
	"io"
)

// AESEncryptCBC_HMAC encrypts the given data using AES-CBC mode with HMAC authentication.
//
// This function provides authenticated encryption by combining AES-CBC encryption
// with HMAC-SHA256 authentication using the Encrypt-then-MAC approach.
//
// Output Format: IV (16 bytes) || Ciphertext || HMAC-SHA256 Tag (32 bytes)
//
// Parameters:
//
//	data   []byte - The plaintext data to be encrypted.
//	encKey []byte - The AES encryption key (must be 16, 24, or 32 bytes for AES-128/192/256).
//	macKey []byte - The HMAC-SHA256 key for authentication. If empty, encKey will be used.
//
// Returns:
//
//	[]byte - The encrypted output: IV (16 bytes) + Ciphertext + HMAC Tag (32 bytes).
//	error  - An error if encryption fails or the key length is invalid.
//
// Technical Details:
//   - IV: 16 bytes (128 bits) - randomly generated for each encryption, equals AES block size
//   - Padding: PKCS#5 padding applied to plaintext before encryption
//   - HMAC: SHA256-based, 32 bytes (256 bits) tag over IV + Ciphertext
//   - Total overhead: 48 bytes (16-byte IV + 32-byte HMAC tag)
//   - Ciphertext length: Padded to multiple of 16 bytes
func AESEncryptCBC_HMAC(data, encKey []byte, macKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	if len(data)%blockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)

	// HMAC(iv || ciphertext)
	if len(macKey) == 0 {
		macKey = encKey
	}
	m := hmac.New(sha256.New, macKey)
	m.Write(iv)
	m.Write(crypted)
	tag := m.Sum(nil)

	out := make([]byte, len(iv)+len(crypted)+len(tag))
	copy(out, iv)
	copy(out[len(iv):], crypted)
	copy(out[len(iv)+len(crypted):], tag)
	return out, nil
}

// AESDecryptCBC_HMAC decrypts the given ciphertext using AES-CBC mode with HMAC verification.
//
// This function performs authenticated decryption by first verifying the HMAC-SHA256
// authentication tag, then decrypting the ciphertext using AES-CBC mode.
//
// Expected Format: IV (16 bytes) || Ciphertext || HMAC-SHA256 Tag (32 bytes)
//
// Parameters:
//
//	crypted []byte - The encrypted data with IV prepended and HMAC tag appended.
//	                 Minimum length: 48 bytes (16-byte IV + 32-byte HMAC tag).
//	encKey  []byte - The AES decryption key (must be 16, 24, or 32 bytes for AES-128/192/256).
//	macKey  []byte - The HMAC-SHA256 key used during encryption. If empty, encKey will be used.
//
// Returns:
//
//	[]byte - The decrypted plaintext data (PKCS#5 padding removed).
//	error  - An error if decryption fails, authentication fails, or the key length is invalid.
//
// Technical Details:
//   - IV: First 16 bytes of crypted data
//   - HMAC Tag: Last 32 bytes of crypted data (verified using constant-time comparison)
//   - Ciphertext: Middle portion, must be multiple of 16 bytes
//   - Minimum input size: 48 bytes
//   - Authentication is verified BEFORE decryption (secure against padding oracle attacks)
func AESDecryptCBC_HMAC(crypted, encKey []byte, macKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	// 🔧 修复：添加 macKey 默认值设置
	if len(macKey) == 0 {
		macKey = encKey
	}

	blockSize := block.BlockSize()
	if len(crypted) < blockSize+sha256.Size || (len(crypted)-blockSize-sha256.Size)%blockSize != 0 {
		return nil, errors.New("invalid ciphertext: insufficient length or invalid padding")
	}
	iv := crypted[:blockSize]
	ct := crypted[blockSize : len(crypted)-sha256.Size]
	tag := crypted[len(crypted)-sha256.Size:]

	// Verify HMAC first (constant time)
	m := hmac.New(sha256.New, macKey) // 现在使用正确的 macKey
	m.Write(iv)
	m.Write(ct)
	expected := m.Sum(nil)
	if subtle.ConstantTimeCompare(expected, tag) != 1 {
		return nil, errors.New("decryption failed")
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	data := make([]byte, len(ct))
	blockMode.CryptBlocks(data, ct)
	data, err = PKCS5UnPadding(data)
	return data, err
}

// AESEncryptCBC_HMAC_PEM encrypts the given data using AES-CBC with HMAC and a PEM-encoded key.
//
// This function decodes a PEM-formatted AES key and uses it for authenticated encryption.
// The PEM block should contain the raw AES key bytes.
//
// Output Format: IV (16 bytes) || Ciphertext || HMAC-SHA256 Tag (32 bytes)
//
// Parameters:
//
//	data   []byte - The plaintext data to be encrypted.
//	key    []byte - The PEM-encoded AES encryption key (containing 16/24/32 raw bytes).
//	macKey []byte - The HMAC key for authentication. If empty, the decoded key will be used.
//
// Returns:
//
//	[]byte - The encrypted output: IV (16 bytes) + Ciphertext + HMAC Tag (32 bytes).
//	error  - An error if encryption fails, PEM decoding fails, or the key is invalid.
//
// Technical Details:
//   - Total overhead: 48 bytes (16-byte IV + 32-byte HMAC tag)
//   - PEM block type can be any valid type containing raw AES key bytes
func AESEncryptCBC_HMAC_PEM(data, key []byte, macKey []byte) ([]byte, error) {
	aeskey, rest := pem.Decode(key)
	if len(rest) > 0 || aeskey == nil {
		return nil, errors.New("invalid aes key")
	}
	return AESEncryptCBC_HMAC(data, aeskey.Bytes, macKey)
}

// AESDecryptCBC_HMAC_PEM decrypts the given ciphertext using AES-CBC with HMAC and a PEM-encoded key.
//
// This function decodes a PEM-formatted AES key and uses it for authenticated decryption.
// The PEM block should contain the raw AES key bytes.
//
// Expected Format: IV (16 bytes) || Ciphertext || HMAC-SHA256 Tag (32 bytes)
//
// Parameters:
//
//	data   []byte - The encrypted data with IV and HMAC tag (minimum 48 bytes).
//	key    []byte - The PEM-encoded AES decryption key (containing 16/24/32 raw bytes).
//	macKey []byte - The HMAC key used during encryption. If empty, the decoded key will be used.
//
// Returns:
//
//	[]byte - The decrypted plaintext data.
//	error  - An error if decryption fails, PEM decoding fails, authentication fails, or the key is invalid.
//
// Technical Details:
//   - Minimum input size: 48 bytes
//   - Authentication is verified before decryption
func AESDecryptCBC_HMAC_PEM(data, key []byte, macKey []byte) ([]byte, error) {
	aeskey, rest := pem.Decode(key)
	if len(rest) > 0 || aeskey == nil {
		return nil, errors.New("invalid aes key")
	}
	return AESDecryptCBC_HMAC(data, aeskey.Bytes, macKey)
}
