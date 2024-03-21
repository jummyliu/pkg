package cryptoutil

import (
	"crypto/elliptic"
	"testing"
)

func TestECCKey(t *testing.T) {
	priKey, pubKey, err := GeneratePEMECCKey(elliptic.P256())
	if err != nil {
		t.Fatalf("Generate ecc key failure: %s", err)
	}
	// pri, _ := x509.ParseECPrivateKey(priKey)
	// pubI, err := x509.ParsePKIXPublicKey(pubKey)
	// pub := pubI.(*ecdsa.PublicKey)
	t.Logf("ECC private key:\n%s", string(priKey))
	t.Logf("ECC public key:\n%s", string(pubKey))
}

func TestECCSign(t *testing.T) {
	priKey := `-----BEGIN ECC PRIVATE KEY-----
MHcCAQEEIBwwFngxcx0ud9qB8isUgj8bUJAfafIeXJy836oqdUXjoAoGCCqGSM49
AwEHoUQDQgAEnpTuRy7GR2bfLZAtd1r9V3Br9AQYpCRhrs6p6lrARWTsJtf6vxR0
FE/T2z3uBtVk/hfOPt5NJM7N12rq0L3PXA==
-----END ECC PRIVATE KEY-----
`
	data := "hello world"
	signature, err := ECCSignPEM([]byte(data), []byte(priKey))
	if err != nil {
		t.Fatalf("data sign failure: %s", err)
	}
	t.Logf("signature:%s", signature)
}

func TestECCVerify(t *testing.T) {
	pubKey := `-----BEGIN ECC PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEnpTuRy7GR2bfLZAtd1r9V3Br9AQY
pCRhrs6p6lrARWTsJtf6vxR0FE/T2z3uBtVk/hfOPt5NJM7N12rq0L3PXA==
-----END ECC PUBLIC KEY-----
`
	data := "hello world"
	signature := "MEYCIQDbZYDtWVDBOsFxGg4Oy3RmIEAIEOKAciIoCr+cqyoWNwIhAOIN4CuVOKRiFDavH19u5DncTDLje0AQsOCRtjm6HQCr"
	err := ECCVerifyPEM([]byte(data), signature, []byte(pubKey))
	if err != nil {
		t.Fatalf("data verify failure: %s", err)
	}
	t.Logf("verify sign success")
}
