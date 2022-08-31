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
	if len(unknowPri) > 0 {
		t.Fatalf("Decode pem pri key failure: %#v", unknowPri)
	}
	if len(unknowPub) > 0 {
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
MIIBCgKCAQEA5ETXgjkEot8li258tpaYFHRAHqpk9TMdn5gJ9ttOTeucBmESj1oQ
Om7K5QGG2cn8zKvALOWDBiQLXlW0N6yNKqEhQApT9t4bBpJyx6jWqdvYj00RT1Bg
joncox9gpFN1bmg5VNQI7rvTqzqJnJmyQwb83OX7uV4JG4ljS8quRg2Xba5ovPgk
Ne+cxCgVwLZt4uC61F8HPsLO+Ck4nenezemiWDlxBxo16zLdY34LMOxr1pcBY0Ah
2jZe4iJ3lFCf8J2EDowMQO5yOnPpUafea9sJnnoRZC9pSUK1FKnY2lpD10VSJj40
DEKhrG5Ef9D5SNDACTK1xGirB5qHBt8vxwIDAQAB
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
		"agent_version": "vvvv1",
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
MIIEpAIBAAKCAQEAzzLf6nbPbyl4GFtLIn/lDLB+5sPf1RdMq770Jh/TywcYiEfI
Cxl7K8n0UdYueKA/lL4Sl6u7c+ovzNu+VePVCScuOMeJKTtl4YZgKrr1Q/DcIcn+
+G5qTx6nRXGWSYPU+El2X0eFnnY+tDMA4Pjh9Zsj9c1XHNdP79U38BTdUhZLZykC
+sOFHfOFm12sUyYdvee+UjGGLhaF88IArD9l/57jZXYgvG23lzIhsQbbNY68rS+8
lKEFEnwAfZJW23uK6peq698+gt4Iyu3DoIfKVbYsy3LnaNYEp99fXqFnpATB99Or
yiOESDp5I4VUpPovW08RG0uqIyXDUsXLzOzXdwIDAQABAoIBAF29dJSNIBf4uP67
/ugi2Ku/Bgq9P30Wx3dEMw00hQtrsGunnhd3dcQ/4CtOLUQhq+JNg/femDY/E1Up
bZlKNE2pzj2d+K4Q020O8F8kqmYMiGs5CgO3YJ1fDupaE1Y7MJUMF2PP5eHlOUUp
dhJSR1ho0gIY3nqL5vxoiKAzsFW8M50MKxs+ETGZ1I2D6o52GinUuLYgRkutte5F
sqaC+ul/DEL7NS+BzEn2zXHjFf0YeLA/FaobtNvmMPeay/uqYP4LGDwQXTUdGUSl
lyy2WM3Z7S07dFpPgkRccZTgh3ynpnkWFzcuJ5lFnUiWUWdcpGigSomQIKiNsyBZ
EMGbJUECgYEA9VhFrEJhALKUt2N5ImDJYy9EkDud/8mQdXSUUAHcBCkUAdj6P+vM
G+Mu/cAHHcjU1Zg7Oo4Pv/Aoc4VVe241GUUqcgWL/V0f/A/eVOGU/4Aa9FPvdt65
UaimwEw9zjEEKZ+dkrVgoXt2xg7/cErZvi1mMOuP3yD/6Ere1PTauecCgYEA2DJ/
LqnM5wX2igymN7ID9xUhmGaWmFWsTKi6kqJuSIlALmRFYOddawssH22wA5nPLcnk
e94uyBTqz95AUHdbTFN9Ftqua1IZ/8yQEH5r6kcnvV6f5qw8MPgVU2ikdJJRH7ho
PfRw1p0PNVCM65KphLuxNeOyPeeWHlMIv5j04/ECgYBFyNrgeWz/9suoMgoVhjQi
GyLEZ8C0LdACKKu66hx7rnd7Yw0jO12uHPuTv5gGl8Y6Dvfh2uCN9rB601USK7G8
w1ikYAGGioN7fcP+nr9zwStpjapSRF2v5Wmwzr7RtE17zWPTg/W9WNHa2g88EH5I
wr8LcSVWERvZJdql9hN0xwKBgQCebysK7D6PkqwgcLKioB8Nw/uRrqRv0GDq8L+B
U+2T1JknJi49nG+2UUKtaXmSufW87XY2XBVWZRXK7WmeTkmmvowt4mXtmgYZkjSF
EdBNqIVz3lM5/UBC9prSPB5AmzU+FKq3tFm4vPJ3NKeAv0LhVZbBEjL98KfvYxRH
LHVSUQKBgQC0RzV2RxBXifBGMFVokCkDi44cdb1Jda8iD69L3vf2BoN0DOwoKmTt
oY4NrzkqOr/siZwAaBuT/z+m9ecaZpcvqvO9lzjM94/4oxXqG3lFffDSxvqJurOR
h9Q2YFOPCRJ/bOd7L/xNvaFY5gfqQZI0gg96L1jBYqDzYFogu8d0dQ==
-----END RSA PRIVATE KEY-----
`
	crypted := "6029e29c82ab414bcced70c04715405fd22a36909e3f4eaf1056214740a7d426106e8cac5d92418ca1ff1121badaf426cfd319a5e207256dc906606f6ba54932fb1f6381018ca28a216fd349d9b069f75a8b97cac3302001621c76f8fd5e30444b1c0d9577f03ee27a4eb7198258ba68f145566a77f7313c7b345dd2f3f794dfaf954897eff794219d4da6f791be5850fc79e1a651778f89cdb50a20c76fe0067c070e74706b8be803b138d227d8ea596bc1f7758b039ebd0e9b9fb4be6deac04eff5e938061c886bf932db203eeeac0311ddc003aaecb9f4a96feb5766adbc7c2f79269a949df77d1ad95b8b1459c289eb0b153d9475ee5153a3b77d023c8e1"
	data, _ := hex.DecodeString(crypted)
	data, err := RSADecryptPEM(data, []byte(priKey))
	if err != nil {
		t.Fatalf("data decode failure: %s", err)
	}
	m := map[string]any{}
	json.Unmarshal(data, &m)
	t.Logf("%#v", m)
}
