package cryptoutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// TODO: 块加密

// GenerateRSAKey 生成 rsa 密钥，bits 可以给 1024
func GenerateRSAKey(bits int) (priKey []byte, pubKey []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	priKey = x509.MarshalPKCS1PrivateKey(privateKey)
	pubKey = x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	return
}

// GeneratePEMRSAKey 生成 PEM格式 的 rsa 密钥，bits 可以给 1024
func GeneratePEMRSAKey(bits int) (priKey []byte, pubKey []byte, err error) {
	pri, pub, err := GenerateRSAKey(bits)
	if err != nil {
		return nil, nil, err
	}
	priKey = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pri,
	})
	pubKey = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pub,
	})
	return priKey, pubKey, nil
}

// RSAEncrypt_ 使用 RSA 公钥加密，接收 *rsa.PublicKey 公钥
// 	OAEP: sha256
func RSAEncrypt_(data []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, data, nil)
}

// RSADecrypt_ 使用 RSA 私钥解密，接收 *rsa.PrivateKey 私钥
// 	OAEP: sha256
func RSADecrypt_(crypted []byte, priKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priKey, crypted, nil)
}

// RSAEncrypt 使用 RSA 公钥加密，接收 []byte 类型的公钥
// 	OAEP: sha256
func RSAEncrypt(data, pubKey []byte) ([]byte, error) {
	pub, err := x509.ParsePKCS1PublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return RSAEncrypt_(data, pub)
}

// RSADecrypt 使用 RSA 私钥解密，接收 []byte 类型的私钥
// 	OAEP: sha256
func RSADecrypt(data, priKey []byte) ([]byte, error) {
	pri, err := x509.ParsePKCS1PrivateKey(priKey)
	if err != nil {
		return nil, err
	}
	return RSADecrypt_(data, pri)
}

// RSAEncryptPEM 使用 RSA 公钥加密，接收PEM格式的公钥
// 	OAEP: sha256
func RSAEncryptPEM(data, pubKey []byte) ([]byte, error) {
	key, rest := pem.Decode(pubKey)
	if len(rest) > 0 {
		return nil, errors.New("invalid public key")
	}
	return RSAEncrypt(data, key.Bytes)
}

// RSADecryptPEM 使用 RSA 私钥解密，接收PEM格式的私钥
// 	OAEP: sha256
func RSADecryptPEM(data, priKey []byte) ([]byte, error) {
	key, rest := pem.Decode(priKey)
	if len(rest) > 0 {
		return nil, errors.New("invalid private key")
	}
	return RSADecrypt(data, key.Bytes)
}
