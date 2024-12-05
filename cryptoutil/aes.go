package cryptoutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/jummyliu/pkg/utils"
)

// GenerateAESKey 生成 256 位 aes 密钥；用 md5 算的 32 个字节
func GenerateAESKey() (key []byte) {
	return []byte(fmt.Sprintf("%x", md5.Sum([]byte(utils.UUID()))))
}

// GeneratePEMAESKey 生成 256 位 PME格式 的 aes 密钥；用 md5 算的 32 个字节
func GeneratePEMAESKey() (key []byte) {
	aesKey := GenerateAESKey()
	key = pem.EncodeToMemory(&pem.Block{
		Type:  "AES KEY",
		Bytes: aesKey,
	})
	return key
}

func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

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

func AESEncryptPEM(data, key []byte) ([]byte, error) {
	aeskey, rest := pem.Decode(key)
	if len(rest) > 0 || aeskey == nil {
		return nil, errors.New("invalid aes key")
	}
	return AESEncrypt(data, aeskey.Bytes)
}

func AESDecryptPEM(data, key []byte) ([]byte, error) {
	aeskey, rest := pem.Decode(key)
	if len(rest) > 0 || aeskey == nil {
		return nil, errors.New("invalid aes key")
	}
	return AESDecrypt(data, aeskey.Bytes)
}
