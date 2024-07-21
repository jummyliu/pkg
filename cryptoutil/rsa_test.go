package cryptoutil

import (
	"bytes"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"testing"
)

func TestRSAKey(t *testing.T) {
	priKey, pubKey, err := GenerateRSAKey(1024)
	if err != nil {
		t.Fatalf("Generate rsa key failure: %s", err)
	}
	pri, _ := x509.ParsePKCS1PrivateKey(priKey)
	pub, _ := x509.ParsePKCS1PublicKey(pubKey)
	t.Log(len(priKey), len(pubKey))
	t.Log(pri.Size(), pri.PublicKey.Size(), pub.Size())
}

func TestRSACrypt(t *testing.T) {
	data := "hello world"
	priKey, pubKey, err := GenerateRSAKey(1024)
	if err != nil {
		t.Fatalf("Generate rsa key failure: %s", err)
	}
	crypted, err := RSAEncrypt([]byte(data), pubKey)
	if err != nil {
		t.Fatalf("RSA encrypt failure: %s", err)
	}
	result, err := RSADecrypt(crypted, priKey)
	if err != nil {
		t.Fatalf("RSA decrypt failure: %s", err)
	}
	if data != string(result) {
		t.Fatalf("RSA %s encrypt -> decrypt %s, not equal!!!", data, string(result))
	}
}

func TestSaveRSA(t *testing.T) {
	priKey, pubKey, _ := GenerateRSAKey(2048)
	priData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: priKey,
	})
	t.Log(string(priData))
	pubData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKey,
	})
	t.Log(string(pubData))

	newPriKey, unknowPri := pem.Decode(priData)
	newPubKey, unknowPub := pem.Decode(pubData)
	if len(unknowPri) > 0 || newPriKey == nil {
		t.Fatalf("Decode pem pri key failure: %#v", unknowPri)
	}
	if len(unknowPub) > 0 || newPubKey == nil {
		t.Fatalf("Decode pem pub key failure: %#v", unknowPub)
	}

	if !bytes.Equal(priKey, newPriKey.Bytes) {
		t.Fatalf("prikey != new prikey: %#v, %#v", priKey, newPriKey)
	}

	if !bytes.Equal(pubKey, newPubKey.Bytes) {
		t.Fatalf("pubkey != new pubkey: %#v, %#v", pubKey, newPubKey)
	}
}

func TestRSAEncode(t *testing.T) {
	pubKey := `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAryX6bELihivaq54fWF3oBgfPg0Rh2zeCU9bmyVYS5aOm6XInaTAa
xJ7h9jZeTmEEuLv9DNTDx7gBNZWfcF8npDg67IBQuelDFLDx6nHV3WbsP2Vslf6y
4ch1LM2GFajjDqOVXJDTtBCnvK43KktLHT5HyGcjMFiRmRzs5i9ojQTXAuP20pkc
VF9BlnnSBEBHQmRl1dFfrjQgWoqlVbScIPgjDJ39FIZxTsZZdqFoPNABbVSuLz0r
LPwqN5Y5DxDvcRO3ewnEsJXPAjeG8GaUONFJt6cvuqXcK4OQc2BeymnhsiQ50vT2
AAQZaLUOV+zp21mKCMZceKzCGP5QuynAMwIDAQAB
-----END RSA PUBLIC KEY-----
`
	data := map[string]any{}
	d, _ := json.Marshal(data)
	crypted, err := RSAEncryptPEM(d, []byte(pubKey))
	if err != nil {
		t.Fatalf("data encode faliure: %s", err)
	}
	t.Logf("%x", crypted)

	data = map[string]any{
		"rule_code":     "x1x1x1",
		"ti_code":       "x2x2x2",
		"agent_version": "vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvvv1vvvv1vvvv1vvvv1111vvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv11111111111vvvvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvvv1vvvv1vvvv1vvvv1111vvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv11111111111vvvvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvvv1vvvv1vvvv1vvvv1111vvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv1vvvv11111111111vv",
	}
	d, _ = json.Marshal(data)

	crypted, err = RSAEncryptPEM(d, []byte(pubKey))
	if err != nil {
		t.Fatalf("data encode faliure: %s", err)
	}
	t.Logf("%x", crypted)
}

func TestRSADecode(t *testing.T) {
	priKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAryX6bELihivaq54fWF3oBgfPg0Rh2zeCU9bmyVYS5aOm6XIn
aTAaxJ7h9jZeTmEEuLv9DNTDx7gBNZWfcF8npDg67IBQuelDFLDx6nHV3WbsP2Vs
lf6y4ch1LM2GFajjDqOVXJDTtBCnvK43KktLHT5HyGcjMFiRmRzs5i9ojQTXAuP2
0pkcVF9BlnnSBEBHQmRl1dFfrjQgWoqlVbScIPgjDJ39FIZxTsZZdqFoPNABbVSu
Lz0rLPwqN5Y5DxDvcRO3ewnEsJXPAjeG8GaUONFJt6cvuqXcK4OQc2BeymnhsiQ5
0vT2AAQZaLUOV+zp21mKCMZceKzCGP5QuynAMwIDAQABAoIBAH/XnL5A489DW01B
EWgSwzUDpngOBc9Y6QwBJFt5NDniBgcHh7TDpAY4Yn6wmI1lS2j77mzbMDwrFtbh
64q+KdU4JepSjpnkpU4JCcsyZARDB9YOVf/19OPQyZZ2PZS5vWIGDROPsrcQIR8b
mrCIXL9voj2o6opzW3MDJfeuSwYCeT+HbxWHg5plyDD7SvWwMkzVnk3OgROZdGIh
4UkMM7XP06UiMc5H341ubCGUyTz8Th2NdUjs9URGt8TacMcqPde/Zua/Wz6jVAV/
15RjMknOhpSSkijdF4qrvupFfpcXH6ORsXX+yBijX7PdSraAclqMNszGyZTFQwRX
8erZz6ECgYEA07Q8xHGKE3LeCfVnx+GGmVrsylZRo+Z5cvaP4Dgju58kUuwxuVo5
otFFNjJEDG6Ey6eiVtcUAj7dAVshqtfihc19jjDZPSOPMfkf7zLoOACZGMO4iw4+
ikq6352N1mgZouV35VGH6oRqxqMyQj9GF8IUvbkjhWd7hwzuglybNNcCgYEA08up
y120lhbSJC1khOAG+mfHLgzBhYZiE3V04ChYybhJtkituzHYn7XK+1KS3BXY7Pr1
To4SWaqFKUUp93Zr2AWDt1Tm8d28Ygjr8G0N+7w6ljBbRDgKwY6QxMD94FGV+RXr
QnYqM+FC1t/8Q+NVJoKCuETyHfe/Bzwf9Lb9CAUCgYBMTLVqB5HAGLI13KCexYWB
V+fntNyPuc0jxgFsyk72nBC3YjE5oG8NY2cSdWNZJ6vsymoT6khn1shIaNPlgxE9
MCaETM6+3kYJuMPtredL58tFxaSJWYToyq43Uc2A7NvwfcuMdqoJt9fT55WBktRs
U6KuDj/jILzAm8SKb13w2QKBgQCmTzjHbo+Ng+IDcnmKNXiFTNSE/pM/zGRbL1JV
apk93S5UqwFxCxU1ZEU90HttwuISRIY35yvVqSbjX2Iy5ZSNjtb9MPggWKPCv4q1
wozGben7YYFpMjCQCOj49yrj6GzBqUqRZ8R/9JTNshifHnYQxU7sb4dHrPEeN0JI
oSBUGQKBgQCRgAncm+nB75AqI5FG9/LYSQmctPDfvu35TPU5U5vl9mNwrq3ANEd3
WdaWzX1AJRcBpA7NodApi7LS2VydAIOrkUTlVZGRIYcfy0D3rhsN4/3+q0hjjZS8
y34eCNgPMxVv47W78ZPOzl7gKAlwD3+ZY2pLLyXYP+4bjl1Yuj7nBw==
-----END RSA PRIVATE KEY-----
`
	crypted := "3b6dc72b3cdc4bf52916d81bfc5f57655782778f8bf5c50c65ca26e8a5f160b2d9775fd2b8e89f923463f99ed80ec789fff28587550f0950feb4e9d5f9a3b9bc2a72ba07daf96958d4ab866a0722d929273cbfe4d3cfe0486399ee6cb7f325a340fd83ad8180fdbb3c97297d397e6718ea1c5e2e743f8531464873fcbdce7f2e99d3fa3c5cac80a437d98624a540f98d399fc87fd4d3564749f93050d5c8636cc045c202fa1222f436db7bf2db07493b015a994a4b42957e6b1a5076e6acdea92c2b76193396455512c1f2fe36cce497c84a53d31e9920a32b4b659c3b857409ba84d02511e9c540bffe4173002c52b3ab6b39393833a095c02aca53c9f3a09a574843167960c59df1d598100b287c2216ba4898057c6ef975678ae6bcb64190e7876d4200129be81a7ed16c170ea829a14f7ea55608567c06e2bf474dc92daaeaf9631cbb6b582e0ebe9bb45accd22e333541c4f46f4c16fbb78e5902ed2ba0a7c0b42ae0ad6148058dcb0149ba6e53b2b9274c181b6bf55c09cfb67bed4972e3dff463a0e80462ecd7d3983b7b63d78cd1bb5fbc26fa15447934c39baed696c96143bb0e0bd0a127af04ba97313ae499b1f164926f6d7aec76638ccd2a75ef88f9451e549e3431edc9e42ffe12fc27e749e2b9ce85fbcddfaba7b068d7ed89816c95c593421a88537c27fe382cc38fae67e69b6c97909d3e3b90cebd888ec512790ccb7be11017cbc0385ef70fa2e1ff481244617cdd31933e7aaa77e33557388ef6485af889e4bff96486ea1e9efc32d0c5cdac65b5580469235f33113a5302b70f7712041c319551d8319bbda68f1fcbf4078cc42d90067dfdab5bf06ddf113825bc2bd3955dc669d0b1064d2e9460b29e76065a0727f88c5aa58f51d669b3a897409138cb23e142053046fb1453f63b114600378564badd82a2116c9d809e667179098f8ed67765b623de786a3ddedd8e31dc1dfe68bd2f2ec5fc10ccb03cc9b09de01424d942563bca48cd6e79e6ee72ed0356b55418c5ca467ebda8dac21acb22da4fcfb38194167548da97bb1f4f510c4674f3f4fa029f4115172a6d"
	data, _ := hex.DecodeString(crypted)
	data, err := RSADecryptPEM(data, []byte(priKey))
	if err != nil {
		t.Fatalf("data decode failure: %s", err)
	}
	m := map[string]any{}
	json.Unmarshal(data, &m)
	t.Logf("%#v", m)
}

func TestRSASign(t *testing.T) {
	priKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAryX6bELihivaq54fWF3oBgfPg0Rh2zeCU9bmyVYS5aOm6XIn
aTAaxJ7h9jZeTmEEuLv9DNTDx7gBNZWfcF8npDg67IBQuelDFLDx6nHV3WbsP2Vs
lf6y4ch1LM2GFajjDqOVXJDTtBCnvK43KktLHT5HyGcjMFiRmRzs5i9ojQTXAuP2
0pkcVF9BlnnSBEBHQmRl1dFfrjQgWoqlVbScIPgjDJ39FIZxTsZZdqFoPNABbVSu
Lz0rLPwqN5Y5DxDvcRO3ewnEsJXPAjeG8GaUONFJt6cvuqXcK4OQc2BeymnhsiQ5
0vT2AAQZaLUOV+zp21mKCMZceKzCGP5QuynAMwIDAQABAoIBAH/XnL5A489DW01B
EWgSwzUDpngOBc9Y6QwBJFt5NDniBgcHh7TDpAY4Yn6wmI1lS2j77mzbMDwrFtbh
64q+KdU4JepSjpnkpU4JCcsyZARDB9YOVf/19OPQyZZ2PZS5vWIGDROPsrcQIR8b
mrCIXL9voj2o6opzW3MDJfeuSwYCeT+HbxWHg5plyDD7SvWwMkzVnk3OgROZdGIh
4UkMM7XP06UiMc5H341ubCGUyTz8Th2NdUjs9URGt8TacMcqPde/Zua/Wz6jVAV/
15RjMknOhpSSkijdF4qrvupFfpcXH6ORsXX+yBijX7PdSraAclqMNszGyZTFQwRX
8erZz6ECgYEA07Q8xHGKE3LeCfVnx+GGmVrsylZRo+Z5cvaP4Dgju58kUuwxuVo5
otFFNjJEDG6Ey6eiVtcUAj7dAVshqtfihc19jjDZPSOPMfkf7zLoOACZGMO4iw4+
ikq6352N1mgZouV35VGH6oRqxqMyQj9GF8IUvbkjhWd7hwzuglybNNcCgYEA08up
y120lhbSJC1khOAG+mfHLgzBhYZiE3V04ChYybhJtkituzHYn7XK+1KS3BXY7Pr1
To4SWaqFKUUp93Zr2AWDt1Tm8d28Ygjr8G0N+7w6ljBbRDgKwY6QxMD94FGV+RXr
QnYqM+FC1t/8Q+NVJoKCuETyHfe/Bzwf9Lb9CAUCgYBMTLVqB5HAGLI13KCexYWB
V+fntNyPuc0jxgFsyk72nBC3YjE5oG8NY2cSdWNZJ6vsymoT6khn1shIaNPlgxE9
MCaETM6+3kYJuMPtredL58tFxaSJWYToyq43Uc2A7NvwfcuMdqoJt9fT55WBktRs
U6KuDj/jILzAm8SKb13w2QKBgQCmTzjHbo+Ng+IDcnmKNXiFTNSE/pM/zGRbL1JV
apk93S5UqwFxCxU1ZEU90HttwuISRIY35yvVqSbjX2Iy5ZSNjtb9MPggWKPCv4q1
wozGben7YYFpMjCQCOj49yrj6GzBqUqRZ8R/9JTNshifHnYQxU7sb4dHrPEeN0JI
oSBUGQKBgQCRgAncm+nB75AqI5FG9/LYSQmctPDfvu35TPU5U5vl9mNwrq3ANEd3
WdaWzX1AJRcBpA7NodApi7LS2VydAIOrkUTlVZGRIYcfy0D3rhsN4/3+q0hjjZS8
y34eCNgPMxVv47W78ZPOzl7gKAlwD3+ZY2pLLyXYP+4bjl1Yuj7nBw==
-----END RSA PRIVATE KEY-----
`
	data := "hello world"
	signature, err := RSASignPKCS1v15PEM([]byte(data), []byte(priKey))
	if err != nil {
		t.Fatalf("data sign failure: %s", err)
	}
	t.Logf("signature:%s", signature)
}

func TestRSAVerify(t *testing.T) {
	pubKey := `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAryX6bELihivaq54fWF3oBgfPg0Rh2zeCU9bmyVYS5aOm6XInaTAa
xJ7h9jZeTmEEuLv9DNTDx7gBNZWfcF8npDg67IBQuelDFLDx6nHV3WbsP2Vslf6y
4ch1LM2GFajjDqOVXJDTtBCnvK43KktLHT5HyGcjMFiRmRzs5i9ojQTXAuP20pkc
VF9BlnnSBEBHQmRl1dFfrjQgWoqlVbScIPgjDJ39FIZxTsZZdqFoPNABbVSuLz0r
LPwqN5Y5DxDvcRO3ewnEsJXPAjeG8GaUONFJt6cvuqXcK4OQc2BeymnhsiQ50vT2
AAQZaLUOV+zp21mKCMZceKzCGP5QuynAMwIDAQAB
-----END RSA PUBLIC KEY-----
`
	data := "hello world"
	signature := "B/IPEZPElGflFJAp2omLMpY4y76LWhNGq0jF0NO+A1Mmrcpi0J3giF7N0kNXe8NZC7er98rNsqKtPjN5cDQ3slLanEOIIa4bfXyIsZDIuTEQz0XM2W9ZUMoJ7uSxcjtbvQBTG/gYWwg90uOTQCHSJ5XbnJY5W5ezCuwTzbDL7GNON3wWRNfkoRzqfM81uRoJMbBQfYSgUF5joROZVAJYeVQRSzBhs/YsMlh68XUlkd5pHyvy+4rAbXvilScfVdeU5v+vkFHtsSvg5hofdpfLNV19IhpZY/8SUFrRGn1YW991Mdhn7jZsqE+41d70FKNyjSHQIz3EW/hYhJJJBx4swA=="
	err := RSAVerifyPKCS1v15PEM([]byte(data), signature, []byte(pubKey))
	if err != nil {
		t.Fatalf("data verify failure: %s", err)
	}
	t.Logf("verify sign success")
}

func TestRSASignPSS(t *testing.T) {
	priKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAryX6bELihivaq54fWF3oBgfPg0Rh2zeCU9bmyVYS5aOm6XIn
aTAaxJ7h9jZeTmEEuLv9DNTDx7gBNZWfcF8npDg67IBQuelDFLDx6nHV3WbsP2Vs
lf6y4ch1LM2GFajjDqOVXJDTtBCnvK43KktLHT5HyGcjMFiRmRzs5i9ojQTXAuP2
0pkcVF9BlnnSBEBHQmRl1dFfrjQgWoqlVbScIPgjDJ39FIZxTsZZdqFoPNABbVSu
Lz0rLPwqN5Y5DxDvcRO3ewnEsJXPAjeG8GaUONFJt6cvuqXcK4OQc2BeymnhsiQ5
0vT2AAQZaLUOV+zp21mKCMZceKzCGP5QuynAMwIDAQABAoIBAH/XnL5A489DW01B
EWgSwzUDpngOBc9Y6QwBJFt5NDniBgcHh7TDpAY4Yn6wmI1lS2j77mzbMDwrFtbh
64q+KdU4JepSjpnkpU4JCcsyZARDB9YOVf/19OPQyZZ2PZS5vWIGDROPsrcQIR8b
mrCIXL9voj2o6opzW3MDJfeuSwYCeT+HbxWHg5plyDD7SvWwMkzVnk3OgROZdGIh
4UkMM7XP06UiMc5H341ubCGUyTz8Th2NdUjs9URGt8TacMcqPde/Zua/Wz6jVAV/
15RjMknOhpSSkijdF4qrvupFfpcXH6ORsXX+yBijX7PdSraAclqMNszGyZTFQwRX
8erZz6ECgYEA07Q8xHGKE3LeCfVnx+GGmVrsylZRo+Z5cvaP4Dgju58kUuwxuVo5
otFFNjJEDG6Ey6eiVtcUAj7dAVshqtfihc19jjDZPSOPMfkf7zLoOACZGMO4iw4+
ikq6352N1mgZouV35VGH6oRqxqMyQj9GF8IUvbkjhWd7hwzuglybNNcCgYEA08up
y120lhbSJC1khOAG+mfHLgzBhYZiE3V04ChYybhJtkituzHYn7XK+1KS3BXY7Pr1
To4SWaqFKUUp93Zr2AWDt1Tm8d28Ygjr8G0N+7w6ljBbRDgKwY6QxMD94FGV+RXr
QnYqM+FC1t/8Q+NVJoKCuETyHfe/Bzwf9Lb9CAUCgYBMTLVqB5HAGLI13KCexYWB
V+fntNyPuc0jxgFsyk72nBC3YjE5oG8NY2cSdWNZJ6vsymoT6khn1shIaNPlgxE9
MCaETM6+3kYJuMPtredL58tFxaSJWYToyq43Uc2A7NvwfcuMdqoJt9fT55WBktRs
U6KuDj/jILzAm8SKb13w2QKBgQCmTzjHbo+Ng+IDcnmKNXiFTNSE/pM/zGRbL1JV
apk93S5UqwFxCxU1ZEU90HttwuISRIY35yvVqSbjX2Iy5ZSNjtb9MPggWKPCv4q1
wozGben7YYFpMjCQCOj49yrj6GzBqUqRZ8R/9JTNshifHnYQxU7sb4dHrPEeN0JI
oSBUGQKBgQCRgAncm+nB75AqI5FG9/LYSQmctPDfvu35TPU5U5vl9mNwrq3ANEd3
WdaWzX1AJRcBpA7NodApi7LS2VydAIOrkUTlVZGRIYcfy0D3rhsN4/3+q0hjjZS8
y34eCNgPMxVv47W78ZPOzl7gKAlwD3+ZY2pLLyXYP+4bjl1Yuj7nBw==
-----END RSA PRIVATE KEY-----
`
	data := "hello world"
	signature, err := RSASignPSSPEM([]byte(data), []byte(priKey))
	if err != nil {
		t.Fatalf("data sign failure: %s", err)
	}
	t.Logf("signature:%s", signature)
}

func TestRSAVerifyPSS(t *testing.T) {
	pubKey := `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAryX6bELihivaq54fWF3oBgfPg0Rh2zeCU9bmyVYS5aOm6XInaTAa
xJ7h9jZeTmEEuLv9DNTDx7gBNZWfcF8npDg67IBQuelDFLDx6nHV3WbsP2Vslf6y
4ch1LM2GFajjDqOVXJDTtBCnvK43KktLHT5HyGcjMFiRmRzs5i9ojQTXAuP20pkc
VF9BlnnSBEBHQmRl1dFfrjQgWoqlVbScIPgjDJ39FIZxTsZZdqFoPNABbVSuLz0r
LPwqN5Y5DxDvcRO3ewnEsJXPAjeG8GaUONFJt6cvuqXcK4OQc2BeymnhsiQ50vT2
AAQZaLUOV+zp21mKCMZceKzCGP5QuynAMwIDAQAB
-----END RSA PUBLIC KEY-----
`
	data := "hello world"
	signature := "Xz17qPx7BYtP80vSxC3ucYrpNAhhj0r5eiJL9hqMPmTgWxGAjuJTHeqSu5r4N7cuIhDWIh1KwxIsk6i1ca+rIG8Qf4EcYYd6gk30qH+q7+Oph9iP7gnISxS2pCcaz+bLo6hvkJos5oNAA7g0+Gz5RZi6wkVI24RZWmk+eQfbpZIkKxXx7pWod+PQiwL5jKDQj4SZrvqkeJhbBUiNoKoLJS14B+zRzej4f3+g2OkP9n6TNJk98gkpljEYN6mPXv78CLogRQe5xP7wHo93fy9Q2fM1fiNajiaSNCLlwYmw+pDvPbua7qv3glGE2VEpYcNXP3mXIDtlOO7kJwLADMpSGA=="
	err := RSAVerifyPSSPEM([]byte(data), signature, []byte(pubKey))
	if err != nil {
		t.Fatalf("data verify failure: %s", err)
	}
	t.Logf("verify sign success")
}
