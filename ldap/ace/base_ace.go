package ace

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/jummyliu/pkg/ldap"
)

// BaseAce
//
//	ACCESS_ALLOWED_ACE, ACCESS_DENIED_ACE, SYSTEM_AUDIT_ACE, SYSTEM_MANDATORY_LABEL_ACE, SYSTEM_SCOPED_POLICY_ID_ACE
//
//	  -----------------------------------------------
//	| 0x00-0x07 | 0x08-0x0F | 0x10-0x17 | 0x18-0x1F |
//	|                    Header                     |
//	|                     Mask                      |
//	|                 Sid(variable)                 |
//	 -----------------------------------------------
type BaseAce struct {
	*ldap.AceHeader
	Mask uint32 `json:"mask"`
	// Mask   [4]byte
	Sid *ldap.Sid `json:"sid"`
}

// SYSTEM_MANDATORY_LABEL_ACE
const (
	SYSTEM_MANDATORY_LABEL_NO_WRITE_UP   ldap.AceMask = 0x00000001
	SYSTEM_MANDATORY_LABEL_NO_READ_UP    ldap.AceMask = 0x00000002
	SYSTEM_MANDATORY_LABEL_NO_EXECUTE_UP ldap.AceMask = 0x00000004
)

// SYSTEM_MANDATORY_LABEL_ACE
const (
	SID_UNTRUSTED_INTEGRITY_LEVEL         uint32 = 0x00000000
	SID_LOW_INTEGRITY_LEVEL               uint32 = 0x00001000
	SID_MEDIUM_INTEGRITY_LEVEL            uint32 = 0x00002000
	SID_HIGH_INTEGRITY_LEVEL              uint32 = 0x00003000
	SID_SYSTEM_INTEGRITY_LEVEL            uint32 = 0x00004000
	SID_PROTECTED_PROCESS_INTEGRITY_LEVEL uint32 = 0x00005000
)

func NewBaseAce(aceBytes []byte) (ace ldap.Ace, err error) {
	if len(aceBytes) < 4 {
		return nil, fmt.Errorf("ace must be at least 4 bytes long")
	}
	aceHeader := &ldap.AceHeader{
		AceType:  aceBytes[0],
		AceFlags: aceBytes[1],
		AceSize:  binary.LittleEndian.Uint16(aceBytes[2:4]),
	}
	if len(aceBytes) < int(aceHeader.AceSize) {
		return nil, fmt.Errorf("ace length is unvalid")
	}
	a := &BaseAce{
		AceHeader: aceHeader,
		Mask:      binary.LittleEndian.Uint32(aceBytes[4:8]),
	}
	offset := 8
	sid, err := ldap.NewSid(aceBytes[offset:a.AceHeader.AceSize])
	if err != nil {
		return nil, fmt.Errorf("failed to parse sid in ace: %s", err)
	}
	a.Sid = sid
	return a, nil
}
func (a *BaseAce) Size() int {
	if a == nil {
		return 0
	}
	return int(a.AceHeader.AceSize)
}

func (a *BaseAce) String() string {
	if a == nil {
		return ""
	}
	buf := strings.Builder{}
	buf.WriteString("        AceType: ")
	buf.WriteString(ldap.AceTypeToFullStr(ldap.AceType(a.AceHeader.AceType)))
	buf.WriteString("\n")
	buf.WriteString("          ")
	buf.WriteString(ldap.AllowedOrDenied(ldap.AceType(a.AceHeader.AceType)))
	buf.WriteString("\n")
	buf.WriteString("        HeaderAceFlags: ")
	buf.WriteString(fmt.Sprintf("0x%02x", a.AceHeader.AceFlags))
	buf.WriteString("\n")
	buf.WriteString("          ")
	buf.WriteString(ldap.AceFlagsToFullStr(ldap.AceFlags(a.AceHeader.AceFlags)))
	buf.WriteString("\n")
	buf.WriteString("        AceSize: ")
	buf.WriteString(fmt.Sprintf("%d", a.AceHeader.AceSize))
	buf.WriteString("\n")
	buf.WriteString("        Mask(Rights): ")
	buf.WriteString(fmt.Sprintf("0x%08x", a.Mask))
	buf.WriteString("\n")
	buf.WriteString("          ")
	buf.WriteString(ldap.AceMaskToFullStr(ldap.AceMask(a.Mask)))
	buf.WriteString("\n")
	buf.WriteString("        Sid: ")
	buf.WriteString(a.Sid.String())
	buf.WriteString("\n")
	return buf.String()
}

// NtString (ace_type;ace_flags;rights;object_guid;inherit_object_guid;account_sid)
func (a *BaseAce) NtString() string {
	if a == nil {
		return ""
	}
	buf := strings.Builder{}
	buf.WriteByte('(')
	buf.WriteString(ldap.AceTypeToStr(ldap.AceType(a.AceHeader.AceType)))
	buf.WriteByte(';')
	buf.WriteString(ldap.AceFlagsToStr(ldap.AceFlags(a.AceHeader.AceFlags)))
	buf.WriteByte(';')
	buf.WriteString(ldap.AceMaskToStr(ldap.AceMask(a.Mask)))
	buf.WriteByte(';')
	// object_guid
	buf.WriteByte(';')
	// inherit_object_guid
	buf.WriteByte(';')
	// sid
	buf.WriteString(a.Sid.Alias())
	buf.WriteByte(')')
	return buf.String()
}

func init() {
	arr := []ldap.AceType{
		ldap.ACCESS_ALLOWED_ACE_TYPE,
		ldap.ACCESS_DENIED_ACE_TYPE,
		ldap.SYSTEM_AUDIT_ACE_TYPE,
		ldap.SYSTEM_MANDATORY_LABEL_ACE_TYPE,
		ldap.SYSTEM_SCOPED_POLICY_ID_ACE_TYPE,
	}
	for _, item := range arr {
		ldap.RegisterAce(item, NewBaseAce)
	}
}
