package cryptoutil

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// RSAEncrypt_ 使用 RSA 公钥加密，接收 *rsa.PublicKey 公钥
//
//	PKCS1v15
func RSAEncryptPKCS1v15_(data []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	buffer := bytes.Buffer{}
	keySize, size := pubKey.Size(), len(data)
	offset := 0
	for offset < size {
		// PKCS1v15: 原文数据长度 <= RSA公钥模长 - 11字节
		endIndex := offset + keySize - 11
		if endIndex > size {
			endIndex = size
		}
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data[offset:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offset = endIndex
	}
	return buffer.Bytes(), nil
}

// RSADecrypt_ 使用 RSA 私钥解密，接收 *rsa.PrivateKey 私钥
//
//	PKCS1v15
func RSADecryptPKCS1v15_(crypted []byte, priKey *rsa.PrivateKey) ([]byte, error) {
	buffer := bytes.Buffer{}
	keySize, size := priKey.Size(), len(crypted)
	offset := 0
	for offset < size {
		endIndex := offset + keySize
		if endIndex > size {
			endIndex = size
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, crypted[offset:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offset = endIndex
	}
	return buffer.Bytes(), nil
}

// RSAEncrypt 使用 RSA 公钥加密，接收 []byte 类型的公钥
//
//	PKCS1v15
func RSAEncryptPKCS1v15(data, pubKey []byte) ([]byte, error) {
	pub, err := x509.ParsePKCS1PublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return RSAEncryptPKCS1v15_(data, pub)
}

// RSADecrypt 使用 RSA 私钥解密，接收 []byte 类型的私钥
//
//	PKCS1v15
func RSADecryptPKCS1v15(data, priKey []byte) ([]byte, error) {
	pri, err := x509.ParsePKCS1PrivateKey(priKey)
	if err != nil {
		return nil, err
	}
	return RSADecryptPKCS1v15_(data, pri)
}

// RSAEncryptPEM 使用 RSA 公钥加密，接收PEM格式的公钥
//
//	PKCS1v15
func RSAEncryptPKCS1v15PEM(data, pubKey []byte) ([]byte, error) {
	key, rest := pem.Decode(pubKey)
	if len(rest) > 0 {
		return nil, errors.New("invalid public key")
	}
	return RSAEncryptPKCS1v15(data, key.Bytes)
}

// RSADecryptPEM 使用 RSA 私钥解密，接收PEM格式的私钥
//
//	PKCS1v15
func RSADecryptPKCS1v15PEM(data, priKey []byte) ([]byte, error) {
	key, rest := pem.Decode(priKey)
	if len(rest) > 0 {
		return nil, errors.New("invalid private key")
	}
	return RSADecryptPKCS1v15(data, key.Bytes)
}
