package ldap_test

import (
	"encoding/base64"
	"testing"

	"github.com/jummyliu/pkg/ldap"
)

func TestNewGuid(t *testing.T) {
	guidBytesB64 := "nwkMjP7CzE+NQmrYBJhVaw=="
	targetGuid := "8c0c099f-c2fe-4fcc-8d42-6ad80498556b"

	guidBytes, err := base64.StdEncoding.DecodeString(guidBytesB64)
	if err != nil {
		t.Fatalf("guidB64 decode failure: %s", err)
	}
	guid, err := ldap.NewGuid(guidBytes)
	if err != nil {
		t.Fatalf("parse guid failure: %s", err)
	}
	// fmt.Println(guid)
	// fmt.Println(ldap.MarshalGUID(guidBytes))
	if guid.String() != targetGuid {
		t.Fatalf("guidStr need %s but got %s", targetGuid, guid)
	}
}

func TestMarshalGUID(t *testing.T) {
	guidBytesB64 := "nwkMjP7CzE+NQmrYBJhVaw=="
	targetGuid := "8c0c099f-c2fe-4fcc-8d42-6ad80498556b"

	guidBytes, err := base64.StdEncoding.DecodeString(guidBytesB64)
	if err != nil {
		t.Fatalf("guidB64 decode failure: %s", err)
	}
	guid, err := ldap.MarshalGUID(guidBytes)
	if err != nil {
		t.Fatalf("parse guid failure: %s", err)
	}
	// fmt.Println(guid)
	// fmt.Println(ldap.MarshalGUID(guidBytes))
	if guid != targetGuid {
		t.Fatalf("guidStr need %s but got %s", targetGuid, guid)
	}
}
