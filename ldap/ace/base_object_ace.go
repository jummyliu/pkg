package ace

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/jummyliu/pkg/ldap"
)

// BaseObjectAce
//
//	ACCESS_ALLOWED_OBJECT_ACE, ACCESS_DENIED_OBJECT_ACE
//
//	 -----------------------------------------------
//	| 0x00-0x07 | 0x08-0x0F | 0x10-0x17 | 0x18-0x1F |
//	|                    Header                     |
//	|                     Mask                      |
//	|                    Flags                      |
//	|         ObjectType(16 bytes, optional)        |
//	|    InheritedObjectType(16 bytes, optional)    |
//	|                 Sid(variable)                 |
//	 -----------------------------------------------
type BaseObjectAce struct {
	*ldap.AceHeader        // [4]byte
	Mask            uint32 `json:"mask"`
	Flags           uint32 `json:"flags"` // AccessAllowedObjectFlags
	// Mask                [4]byte
	// Flags               [4]byte
	ObjectType          Guid      `json:"objectType"`
	InheritedObjectType Guid      `json:"inheritedObjectType"`
	Sid                 *ldap.Sid `json:"sid"`
}

type Guid [16]byte

func (g Guid) MarshalJSON() ([]byte, error) {
	tmp, err := ldap.NewGuid(g[:])
	if err != nil {
		return []byte("\"\""), nil
	}
	tmpStr := fmt.Sprintf("\"%s\"", tmp.String())
	return []byte(tmpStr), nil
}

// type AccessAllowedObjectFlags uint32

// const (
// 	Null                              AccessAllowedObjectFlags = 0x00000000
// 	ACE_OBJECT_TYPE_PRESENT           AccessAllowedObjectFlags = 0x00000001
// 	ACE_INHERITED_OBJECT_TYPE_PRESENT AccessAllowedObjectFlags = 0x00000002
// )

// var AccessAllowedObjectFlagsMap = map[AccessAllowedObjectFlags]string{
// 	// Null:                              "Null",
// 	ACE_OBJECT_TYPE_PRESENT:           "ACE_OBJECT_TYPE_PRESENT",
// 	ACE_INHERITED_OBJECT_TYPE_PRESENT: "ACE_INHERITED_OBJECT_TYPE_PRESENT",
// }

// func AccessAllowedObjectFlagsToStr(flags AccessAllowedObjectFlags) string {
// 	var tmpArr = []string{}
// 	for key, val := range AccessAllowedObjectFlagsMap {
// 		if key&flags == key {
// 			tmpArr = append(tmpArr, val)
// 		}
// 	}
// 	return strings.Join(tmpArr, " ")
// }
// func AccessAllowedObjectFlagsToFullStr(flags AccessAllowedObjectFlags) string {
// 	var tmpArr = []string{}
// 	for key, val := range AccessAllowedObjectFlagsMap {
// 		if key&flags == key {
// 			tmpArr = append(tmpArr, fmt.Sprintf("%s(0x%08x)", val, key))
// 		}
// 	}
// 	return strings.Join(tmpArr, ",")
// }

// specifies mask
//
// ACCESS_ALLOWED_OBJECT_ACE, ACCESS_DENIED_OBJECT_ACE
// ACCESS_ALLOWED_CALLBACK_OBJECT_ACE, ACCESS_DENIED_CALLBACK_OBJECT_ACE, SYSTEM_AUDIT_OBJECT_ACE, SYSTEM_AUDIT_CALLBACK_OBJECT_ACE
const (
	ADS_RIGHT_DS_CONTROL_ACCESS ldap.AceMask = 0x00000100
	ADS_RIGHT_DS_CREATE_CHILD   ldap.AceMask = 0x00000001
	ADS_RIGHT_DS_DELETE_CHILD   ldap.AceMask = 0x00000002
	ADS_RIGHT_DS_READ_PROP      ldap.AceMask = 0x00000010
	ADS_RIGHT_DS_WRITE_PROP     ldap.AceMask = 0x00000020
	ADS_RIGHT_DS_SELF           ldap.AceMask = 0x00000008
)

func NewBaseObjectAce(aceBytes []byte) (ace ldap.Ace, err error) {
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
	a := &BaseObjectAce{
		AceHeader: aceHeader,
		Mask:      binary.LittleEndian.Uint32(aceBytes[4:8]),
		Flags:     binary.LittleEndian.Uint32(aceBytes[8:12]),
		// ObjectType:          [16]byte(aceBytes[12:28]),
		// InheritedObjectType: [16]byte(aceBytes[28:44]),
	}
	offset := 12
	var objectType [16]byte
	if ldap.ACE_OBJECT_TYPE_PRESENT&ldap.AccessAllowedObjectFlags(a.Flags) == ldap.ACE_OBJECT_TYPE_PRESENT {
		objectType = [16]byte(aceBytes[offset : offset+16])
		offset += 16
	}
	var inheritedObjectType [16]byte
	if ldap.ACE_INHERITED_OBJECT_TYPE_PRESENT&ldap.AccessAllowedObjectFlags(a.Flags) == ldap.ACE_INHERITED_OBJECT_TYPE_PRESENT {
		inheritedObjectType = [16]byte(aceBytes[offset : offset+16])
		offset += 16
	}
	a.ObjectType = objectType
	a.InheritedObjectType = inheritedObjectType
	sid, err := ldap.NewSid(aceBytes[offset:a.AceHeader.AceSize])
	if err != nil {
		return nil, fmt.Errorf("failed to parse sid in ace: %s", err)
	}
	a.Sid = sid
	return a, nil
}

func (a *BaseObjectAce) Size() int {
	if a == nil {
		return 0
	}
	return int(a.AceHeader.AceSize)
}

func (a *BaseObjectAce) String() string {
	if a == nil {
		return ""
	}
	objectType, _ := ldap.NewGuid(a.ObjectType[:])
	inheritedObjectType, _ := ldap.NewGuid(a.InheritedObjectType[:])
	buf := strings.Builder{}
	buf.WriteString("        AceType: ")
	buf.WriteString(ldap.AceTypeToFullStr(ldap.AceType(a.AceHeader.AceType)))
	buf.WriteString("\n")
	buf.WriteString("          ")
	buf.WriteString(ldap.AllowedOrDenied(ldap.AceType(a.AceHeader.AceType)))
	if ldap.IsInherit(ldap.AccessAllowedObjectFlags(a.Flags)) {
		buf.WriteString(" Inherit")
	}
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
	buf.WriteString("        Flags: ")
	buf.WriteString(fmt.Sprintf("0x%08x", a.Flags))
	buf.WriteString("\n")
	buf.WriteString("          ")
	buf.WriteString(ldap.AccessAllowedObjectFlagsToFullStr(ldap.AccessAllowedObjectFlags(a.Flags)))
	buf.WriteString("\n")
	buf.WriteString("        ObjectType: ")
	buf.WriteString(objectType.Alias())
	buf.WriteString("\n")
	buf.WriteString("        InheritedObjectType: ")
	buf.WriteString(inheritedObjectType.Alias())
	buf.WriteString("\n")
	buf.WriteString("        Sid: ")
	buf.WriteString(a.Sid.Alias())
	buf.WriteString("\n")
	return buf.String()
}

// NtString (ace_type;ace_flags;rights;object_guid;inherit_object_guid;account_sid)
func (a *BaseObjectAce) NtString() string {
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
	objectType, _ := ldap.NewGuid(a.ObjectType[:])
	buf.WriteString(objectType.String())
	buf.WriteByte(';')
	// inherit_object_guid
	inheritedObjectType, _ := ldap.NewGuid(a.InheritedObjectType[:])
	buf.WriteString(inheritedObjectType.String())
	buf.WriteByte(';')
	// sid
	buf.WriteString(a.Sid.String())
	buf.WriteByte(')')
	return buf.String()
}

func init() {
	arr := []ldap.AceType{
		ldap.ACCESS_ALLOWED_OBJECT_ACE_TYPE,
		ldap.ACCESS_DENIED_OBJECT_ACE_TYPE,
	}
	for _, item := range arr {
		ldap.RegisterAce(item, NewBaseObjectAce)
	}

}
