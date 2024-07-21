package cryptoutil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// GenerateECCKey 生成 ecc 密钥，c 可以是
//
//	elliptic.P224()
//	elliptic.P256()
//	elliptic.P384()
//	elliptic.P521()
//
// 或其他椭圆曲线算法
func GenerateECCKey(c elliptic.Curve) (priKey []byte, pubKey []byte, err error) {
	privateKey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	priKey, err = x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err = x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return
}

// GeneratePEMECCKey 生成 PEM格式 的 ecc 密钥，c 可以是
//
//	elliptic.P224()
//	elliptic.P256()
//	elliptic.P384()
//	elliptic.P521()
//
// 或其他椭圆曲线算法
func GeneratePEMECCKey(c elliptic.Curve) (priKey []byte, pubKey []byte, err error) {
	pri, pub, err := GenerateECCKey(c)
	if err != nil {
		return nil, nil, err
	}
	priKey = pem.EncodeToMemory(&pem.Block{
		Type:  "ECC PRIVATE KEY",
		Bytes: pri,
	})
	pubKey = pem.EncodeToMemory(&pem.Block{
		Type:  "ECC PUBLIC KEY",
		Bytes: pub,
	})
	return priKey, pubKey, nil
}

// ECCSign_ 使用 ecc 私钥签名，接收 *ecdsa.PrivateKey 公钥
func ECCSign_(data []byte, priKey *ecdsa.PrivateKey) (sigB64 string, err error) {
	hashAlgorithm := crypto.SHA256
	instance := hashAlgorithm.New()
	instance.Write(data)
	hashed := instance.Sum(nil)
	sig, err := ecdsa.SignASN1(rand.Reader, priKey, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

var ErrEccVerification = errors.New("ecdsa: verification error")

// ECCVerify_ 使用 ecc 公钥验签，接收 *ecdsa.PublicKey 公钥
func ECCVerify_(data []byte, sigB64 string, pubKey *ecdsa.PublicKey) error {
	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return err
	}
	hashAlgorithm := crypto.SHA256
	instance := hashAlgorithm.New()
	instance.Write(data)
	hashed := instance.Sum(nil)
	result := ecdsa.VerifyASN1(pubKey, hashed, sig)
	if !result {
		return ErrEccVerification
	}
	return nil
}

// ECCSign 使用 ecc 私钥签名，接收 []byte 类型的公钥
func ECCSign(data, priKey []byte) (string, error) {
	pri, err := x509.ParseECPrivateKey(priKey)
	if err != nil {
		return "", err
	}
	return ECCSign_(data, pri)
}

// ECCVerify 使用 ecc 公钥验签，接收 []byte 类型的公钥
func ECCVerify(data []byte, sigB64 string, pubKey []byte) error {
	pubI, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return err
	}
	if pub, ok := pubI.(*ecdsa.PublicKey); ok {
		return ECCVerify_(data, sigB64, pub)
	}
	return errors.New("invalid ecc public key")
}

// ECCSignPEM 使用 ecc 私钥签名，接收PEM格式的公钥
func ECCSignPEM(data, priKey []byte) (string, error) {
	key, rest := pem.Decode(priKey)
	if len(rest) > 0 || key == nil {
		return "", errors.New("invalid ecc private key")
	}
	return ECCSign(data, key.Bytes)
}

// ECCVerifyPEM 使用 ecc 公钥验签，接收PEM格式的公钥
func ECCVerifyPEM(data []byte, signStr string, pubKey []byte) error {
	key, rest := pem.Decode(pubKey)
	if len(rest) > 0 || key == nil {
		return errors.New("invalid ecc public key")
	}
	return ECCVerify(data, signStr, key.Bytes)
}
