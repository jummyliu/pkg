package ldap_test

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/jummyliu/pkg/ldap"
)

func TestNewSid(t *testing.T) {
	sidBytesB64 := "AQUAAAAAAAUVAAAAyLgUr3zUJld3XHya6AMAAA=="
	targetSid := "S-1-5-21-2937370824-1462162556-2591841399-1000"

	sidBytes, err := base64.StdEncoding.DecodeString(sidBytesB64)
	if err != nil {
		t.Fatalf("sidB64 decode failure: %s", err)
	}
	sid, err := ldap.NewSid(sidBytes)
	if err != nil {
		t.Fatalf("parse sid failure: %s", err)
	}
	// fmt.Println(sid)
	// fmt.Println(ldap.MarshalSid(sidBytes))
	if sid.String() != targetSid {
		t.Fatalf("sidStr need %s but got %s", targetSid, sid)
	}
}

func TestMarshalSID(t *testing.T) {
	sidBytesB64 := "AQUAAAAAAAUVAAAAyLgUr3zUJld3XHya6AMAAA=="
	targetSid := "S-1-5-21-2937370824-1462162556-2591841399-1000"

	sidBytes, err := base64.StdEncoding.DecodeString(sidBytesB64)
	if err != nil {
		t.Fatalf("sidB64 decode failure: %s", err)
	}
	sid, err := ldap.MarshalSid(sidBytes)
	if err != nil {
		t.Fatalf("parse sid failure: %s", err)
	}
	// fmt.Println(sid)
	// fmt.Println(ldap.MarshalSid(sidBytes))
	if sid != targetSid {
		t.Fatalf("sidStr need %s but got %s", targetSid, sid)
	}
}

func TestLittleEndian(t *testing.T) {
	data := []byte{0b01000000, 0b00000001}
	val := binary.LittleEndian.Uint16(data)
	fmt.Println(val, val&64, val&256)
}
