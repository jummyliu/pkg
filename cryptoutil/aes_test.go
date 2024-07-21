package cryptoutil

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"testing"
)

const aesKey = "d48061d020a5bd1589dca4d198104238"

func TestAESKey(t *testing.T) {
	key := GenerateAESKey()
	if len(key) != 32 {
		t.Fatalf("Generate aes key len %d(%d bit), want %d(%d bit)", len(key), len(key)*8, 32, 32*8)
	}
}

func TestAESCrypt(t *testing.T) {
	data := "hello world"
	key := GenerateAESKey()
	crypted, err := AESEncrypt([]byte(data), key)
	if err != nil {
		t.Fatalf("AES encrypt failure: %s", err)
	}
	fmt.Printf("%x", crypted)
	result, err := AESDecrypt(crypted, key)
	if err != nil {
		t.Fatalf("AES decrypt failure: %s", err)
	}
	if data != string(result) {
		t.Fatalf("AES %s encrypt -> decrypt %s, not equal!!!", data, string(result))
	}
}

func TestSaveAES(t *testing.T) {
	key := GenerateAESKey()
	priData := pem.EncodeToMemory(&pem.Block{
		Type:  "AES KEY",
		Bytes: key,
	})
	t.Log(string(priData))

	newKey, unknowKey := pem.Decode(priData)
	if len(unknowKey) > 0 || newKey == nil {
		t.Fatalf("Decode pem unknow key failure: %#v", unknowKey)
	}
	if !bytes.Equal(key, newKey.Bytes) {
		t.Fatalf("key != new key: %#v, %#v", key, newKey)
	}
}

func TestAESEncode(t *testing.T) {
	// key := "d48061d020a5bd1589dca4d198104238"
	data := map[string]any{}
	d, _ := json.Marshal(&data)
	crypted, err := AESEncrypt(d, []byte(aesKey))
	if err != nil {
		t.Fatalf("data encode aes failure: %s", err)
	}
	t.Logf("%x", crypted)

	data = map[string]any{
		"rule_code": "x1x2x3",
	}
	d, _ = json.Marshal(&data)
	crypted, err = AESEncrypt(d, []byte(aesKey))
	if err != nil {
		t.Fatalf("data encode aes failure: %s", err)
	}
	t.Logf("%x", crypted)
}

func TestAESDecode(t *testing.T) {
	// key := "d48061d020a5bd1589dca4d198104238"
	crypted := "00161b21e43be619246913f17004fd54d3f7dd981f1a9f763addc7954ea8c00f50fa2145b40a7df2cbea0c278cc028a77d0f7623170b119d34c57558dbdd98d74ce6347c6b9691c39d08b82c65621625d7782ae984a753da72650b92a8ab2511881a460bd1ed9d4e147a05aa03252741a8802607eb1389752658e2677a25a433f38e570dbe4fc5fe6f398614c3f9bf4cdb38f02a5c07eb2b2d32153507a5b1e3c53d5d350b46a743c5f9c0067a7957a21e2aa624a41220ee524ec7f8314a084abadb50da07978ba9b79c99a2d11e9b11f91534f649c74664a7fda9862a1dd983a5da68b053e3147ede8d42ad86d5625de6320d7abd32a757cc7bc1eba3c22fe9c5ae41e2c907c7e641bd7bac4663d071113aa0d6e9255e3429aed1abca66481ba19ca5d38df6600c5cb945df6257f78169981776bf9a361ce3a3daf86b76b0b3c114fa5c5115a63e333598e40d60c8925ee967e88f5e7d548f9c32d7dbe75b52eda15ed1f28f24c6d060e3b5356bde71d82ce52409f0d59f0b4bc2b1b0b1dbc52d3f9f67958a12deb6fa5f769a3e85c44f9e8a4a9b3d3e40151ed14f534676c7f836dbcbd800b7e8df4157c6fec17700e6c5dcabf66f25ed55dbf72f6df54edb372d83d0f101792e239e455ac6d413ba454cb8f4c4991d3a8c95f3f9a79c06f0cddb4f0d79e138eb07320155059b062cfa3ba75de663a5e72816d05d6565dc02a1b367c82dbbccfa577d3f87fc3f5642f90ed73c0e44f0b554365b2f6845bc6413e41f129e0ef14ea1e75b83a8362c4ee98597c08f1df1fe99205d8ec0e8808a6fe6d57170fd5c6d14fdd98984d2df13bdbc2e4f5c91eb51c73d2e24f94d24a452ac453a992f4f91eb6e524ac219827567b6f14d8ac2bede536c47d7a5a3d90937d9e67e3ea873d1dfdaaddcb5989286591fa3c864efedbb09e8fe3fba2c148b4506d866fa2dba437a2630c9aee7e85463fb58ad6b6d9c5873f131b0b2d15fcbdda053aa61d072658e55517f3f5f3e9d0003a5103adcb85a6c579d336e1c178f05232f4fc992a4ed5932d4b8c47e1ec9bacd0c96ec07a5e4cdff2e00f69eaf1b8fd1901e6463e237d40e93f743c90db0f5fc9b5b891f1dc16432d09a2e78285cfe8d8e5e2d5df88bebf9a54aff0d2a6b09d4229e6932e1e793c268bcb770e30159bf8df5625f98f4fb433144fba9eb56ebb5e752196e170ddc8bd18e408fddc8d665a4e5540b155ffb5b29e4074ef945"
	data, _ := hex.DecodeString(crypted)
	data, err := AESDecrypt(data, []byte(aesKey))
	if err != nil {
		t.Fatalf("data decode aes failure: %s", err)
	}
	m := map[string]any{}
	json.Unmarshal(data, &m)
	t.Logf("%#v", m)
}

func TestHexDump(t *testing.T) {
	str := "adlkhfd的撒开发撒赖扩大解放"

	t.Log(hex.Dump([]byte(str)))
}

func TestAESEncodeRule(t *testing.T) {
	data := map[string]any{
		"rule_code": "e1044cdc-a94e-4d52-bf68-a60d7525d368",
	}
	d, _ := json.Marshal(&data)
	crypted, err := AESEncrypt(d, []byte(aesKey))
	if err != nil {
		t.Fatalf("data encode aes failure: %s", err)
	}
	t.Logf("%x", crypted)
}
func TestAESDecodeRule(t *testing.T) {
	crypted := "00161b21e43be619246913f17004fd54d3f7dd981f1a9f763addc7954ea8c00f50fa2145b40a7df2cbea0c278cc028a77d0f7623170b119d34c57558dbdd98d74ce6347c6b9691c39d08b82c65621625d7782ae984a753da72650b92a8ab2511881a460bd1ed9d4e147a05aa03252741a8802607eb1389752658e2677a25a433f38e570dbe4fc5fe6f398614c3f9bf4cdb38f02a5c07eb2b2d32153507a5b1e3c53d5d350b46a743c5f9c0067a7957a21e2aa624a41220ee524ec7f8314a084abadb50da07978ba9b79c99a2d11e9b11f91534f649c74664a7fda9862a1dd983a5da68b053e3147ede8d42ad86d5625de6320d7abd32a757cc7bc1eba3c22fe9c5ae41e2c907c7e641bd7bac4663d071113aa0d6e9255e3429aed1abca66481ba19ca5d38df6600c5cb945df6257f78169981776bf9a361ce3a3daf86b76b0b3c114fa5c5115a63e333598e40d60c8925ee967e88f5e7d548f9c32d7dbe75b52eda15ed1f28f24c6d060e3b5356bde71d82ce52409f0d59f0b4bc2b1b0b1dbc52d3f9f67958a12deb6fa5f769a3e85c44f9e8a4a9b3d3e40151ed14f534676c7f836dbcbd800b7e8df4157c6fec17700e6c5dcabf66f25ed55dbf72f6df54edb372d83d0f101792e239e455ac6d413ba454cb8f4c4991d3a8c95f3f9a79c06f0cddb4f0d79e138eb07320155059b062cfa3ba75de663a5e72816d05d6565dc02a1b367c82dbbccfa577d3f87fc3f5642f90ed73c0e44f0b554365b2f6845bc6413e41f129e0ef14ea1e75b83a8362c4ee98597c08f1df1fe99205d8ec0e8808a6fe6d57170fd5c6d14fdd98984d2df13bdbc2e4f5c91eb51c73d2e24f94d24a452ac453a992f4f91eb6e524ac219827567b6f14d8ac2bede536c47d7a5a3d90937d9e67e3ea873d1dfdaaddcb5989286591fa3c864efedbb09e8fe3fba2c148b4506d866fa2dba437a2630c9aee7e85463fb58ad6b6d9c5873f131b0b2d15fcbdda053aa61d072658e55517f3f5f3e9d0003a5103adcb85a6c579d336e1c178f05232f4fc992a4ed5932d4b8c47e1ec9bacd0c96ec07a5e4cdff2e00f69eaf1b8fd1901e6463e237d40e93f743c90db0f5fc9b5b891f1dc16432d09a2e78285cfe8d8e5e2d5df88bebf9a54aff0d2a6b09d4229e6932e1e793c268bcb770e30159bf8df5625f98f4fb433144fba9eb56ebb5e752196e170ddc8bd18e408fddc8d665a4e5540b155ffb5b29e4074ef945"
	data, _ := hex.DecodeString(crypted)
	data, err := AESDecrypt(data, []byte(aesKey))
	if err != nil {
		t.Fatalf("data decode aes failure: %s", err)
	}
	m := map[string]any{}
	json.Unmarshal(data, &m)
	t.Logf("%#v", m)
}

func TestAESEncodeTI(t *testing.T) {
	data := map[string]any{
		"to_code":   "5ca5da99-cc8a-4367-86f0-ab37f49049c3",
		"from_code": "",
		"page":      1,
		"size":      10,
	}
	d, _ := json.Marshal(&data)
	crypted, err := AESEncrypt(d, []byte(aesKey))
	if err != nil {
		t.Fatalf("data encode aes failure: %s", err)
	}
	t.Logf("%x", crypted)
}
func TestAESDecodeTI(t *testing.T) {
	crypted := "7c3e29085f4c78637b075eb8d916498093f4d311eae3b6480aaf9ecbb6533e3c0a9399f92868ac8efec3ce1258eaedfd74fc4e8425a7fee92dd2c97ff6f00153c9124ee609beea1220aeabbdab31df475b10f520520923f14fda0410f97efd726661e20a2a82d3069f5291a9dcbc5ac3e9a93f46d6efed1de357a5346cac82ef1c8a07cd2fc07f9c27cdf7e26e31a17fcdf2b646a6c28eb9eb71d38755bb365c68eb2d135d4ec2d15273f490bdf9d090e38034147591e38a002b34eba77bb96ce0ad966cd5702d2a5ba3664d86af5a98231acd36d9401a35efa1eea31e5094b16600ab1e704c4dc82aa3caa041a98778511a28d47f62030814d14587d1414e803272f38bf86939a982471d1add95f96767afdb21c034f4d4a98be390434ac4b9ab96a2581a04e66067b273c47f8c07ddbc54940d98970b89a02ced824a273c82592e8860bdf70773986b9973abb6e970decbc4106c6bc5c586405ab79ab3663293efba6a5754090874aa2d55499417cb0665066322e73d75185c417840e886311bdbeda693b16ce8b4df439f74373274eccc9328fe3b7c16d00d0b747245c4223b2b9a0819607fd3e532393f02a89131a8c35d6377cfaa317ed06876fd68360b"
	data, _ := hex.DecodeString(crypted)
	data, err := AESDecrypt(data, []byte(aesKey))
	if err != nil {
		t.Fatalf("data decode aes failure: %s", err)
	}
	m := map[string]any{}
	json.Unmarshal(data, &m)
	t.Logf("%#v", m)
}
