package ldap

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/jummyliu/pkg/utils"
)

// NtSecurityDescriptor
//
//	 -----------------------------------------------
//	| 0x00-0x07 | 0x08-0x0F | 0x10-0x17 | 0x18-0x1F |
//	| Revision  | Sbz1      |        Control        |
//	|                  OffsetOwner                  |
//	|                  OffsetGroup                  |
//	|                  OffsetSacl                   |
//	|                  OffsetDacl                   |
//	|               OwnerSid(variable)              |
//	|               GroupSid(variable)              |
//	|                 Sacl(variable)                |
//	|                 Dacl(variable)                |
//	 -----------------------------------------------
type NtSecurityDescriptor struct {
	Revision    byte   `json:"-"`
	Sbz1        byte   `json:"-"`
	Control     uint16 `json:"control"` // LittleEndian
	OffsetOwner uint32 `json:"-"`       // LittleEndian
	OffsetGroup uint32 `json:"-"`       // LittleEndian
	OffsetSacl  uint32 `json:"-"`       // LittleEndian
	OffsetDacl  uint32 `json:"-"`       // LittleEndian
	OwnerSid    *Sid   `json:"ownerSid"`
	GroupSid    *Sid   `json:"groupSid"`
	Sacl        *Acl   `json:"sacl"`
	Dacl        *Acl   `json:"dacl"`
}

type SDDLControl uint16

const (
	SE_SELF_RELATIVE         SDDLControl = 0x8000 // SR
	SE_RM_CONTROL_VALID      SDDLControl = 0x4000 // RM
	SE_SACL_PROTECTED        SDDLControl = 0x2000 // PS
	SE_DACL_PROTECTED        SDDLControl = 0x1000 // PD
	SE_SACL_AUTO_INHERITED   SDDLControl = 0x0800 // SI
	SE_DACL_AUTO_INHERITED   SDDLControl = 0x0400 // DI
	SE_SACL_AUTO_INHERIT_REQ SDDLControl = 0x0200 // SC
	SE_DACL_AUTO_INHERIT_REQ SDDLControl = 0x0100 // DC
	SE_SERVER_SECURITY       SDDLControl = 0x0080 // SS
	SE_DACL_TRUSTED          SDDLControl = 0x0040 // DT
	SE_SACL_DEFAULTED        SDDLControl = 0x0020 // SD
	SE_SACL_PRESENT          SDDLControl = 0x0010 // SP
	SE_DACL_DEFAULTED        SDDLControl = 0x0008 // DD
	SE_DACL_PRESENT          SDDLControl = 0x0004 // DP
	SE_GROUP_DEFAULTED       SDDLControl = 0x0002 // GD
	SE_OWNER_DEFAULTED       SDDLControl = 0x0001 // OD
)

var SDDLControlMap = map[SDDLControl]string{
	SE_SELF_RELATIVE:         "SR",
	SE_RM_CONTROL_VALID:      "RM",
	SE_SACL_PROTECTED:        "PS",
	SE_DACL_PROTECTED:        "PD",
	SE_SACL_AUTO_INHERITED:   "SI",
	SE_DACL_AUTO_INHERITED:   "DI",
	SE_SACL_AUTO_INHERIT_REQ: "SC",
	SE_DACL_AUTO_INHERIT_REQ: "DC",
	SE_SERVER_SECURITY:       "SS",
	SE_DACL_TRUSTED:          "DT",
	SE_SACL_DEFAULTED:        "SD",
	SE_SACL_PRESENT:          "SP",
	SE_DACL_DEFAULTED:        "DD",
	SE_DACL_PRESENT:          "DP",
	SE_GROUP_DEFAULTED:       "GD",
	SE_OWNER_DEFAULTED:       "OD",
}

var SDDLControlFullMap = map[SDDLControl]string{
	SE_SELF_RELATIVE:         "SELF_RELATIVE",
	SE_RM_CONTROL_VALID:      "RM_CONTROL_VALID",
	SE_SACL_PROTECTED:        "SACL_PROTECTED",
	SE_DACL_PROTECTED:        "DACL_PROTECTED",
	SE_SACL_AUTO_INHERITED:   "SACL_AUTO_INHERITED",
	SE_DACL_AUTO_INHERITED:   "DACL_AUTO_INHERITED",
	SE_SACL_AUTO_INHERIT_REQ: "SACL_COMPUTED_INHERITANCE_REQUIRED",
	SE_DACL_AUTO_INHERIT_REQ: "DACL_COMPUTED_INHERITANCE_REQUIRED",
	SE_SERVER_SECURITY:       "SERVER_SECURITY",
	SE_DACL_TRUSTED:          "DACL_TRUSTED",
	SE_SACL_DEFAULTED:        "SACL_DEFAULTED",
	SE_SACL_PRESENT:          "SACL_PRESENT",
	SE_DACL_DEFAULTED:        "DACL_DEFAULTED",
	SE_DACL_PRESENT:          "DACL_PRESENT",
	SE_GROUP_DEFAULTED:       "GROUP_DEFAULTED",
	SE_OWNER_DEFAULTED:       "OWNER_DEFAULTED",
}

func SDDLControlToStr(control SDDLControl) string {
	var tmpArr = []string{}
	for key, val := range SDDLControlMap {
		if key&control == key {
			tmpArr = append(tmpArr, val)
		}
	}
	return strings.Join(tmpArr, "")
}

func SDDLControlToFullStr(control SDDLControl) string {
	var tmpArr = []string{}
	for key, val := range SDDLControlFullMap {
		if key&control == key {
			tmpArr = append(tmpArr, fmt.Sprintf("%s(0x%04x)", val, key))
		}
	}
	return strings.Join(tmpArr, ",")
}

func NewNtSecurityDescriptor(descBytes []byte) (descriptor *NtSecurityDescriptor, err error) {
	if len(descBytes) < 20 {
		return nil, fmt.Errorf("ntSecurityDescriptord must be at least 20 bytes long")
	}
	descriptor = &NtSecurityDescriptor{
		Revision:    descBytes[0],
		Sbz1:        descBytes[1],
		Control:     binary.LittleEndian.Uint16(descBytes[2:4]),
		OffsetOwner: binary.LittleEndian.Uint32(descBytes[4:8]),
		OffsetGroup: binary.LittleEndian.Uint32(descBytes[8:12]),
		OffsetSacl:  binary.LittleEndian.Uint32(descBytes[12:16]),
		OffsetDacl:  binary.LittleEndian.Uint32(descBytes[16:20]),
	}
	if descriptor.OffsetOwner != 0 {
		ownerSid, err := NewSid(descBytes[descriptor.OffsetOwner:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse ownerSid bytes: %s", err)
		}
		descriptor.OwnerSid = ownerSid
	}
	if descriptor.OffsetGroup != 0 {
		groupSid, err := NewSid(descBytes[descriptor.OffsetGroup:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse ownerGroup bytes: %s", err)
		}
		descriptor.GroupSid = groupSid
	}
	if SE_SACL_PRESENT&SDDLControl(descriptor.Control) == SE_SACL_PRESENT {
		// parse SACL
		sacl, err := NewAcl(descBytes[descriptor.OffsetSacl:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse sacl bytes: %s", err)
		}
		descriptor.Sacl = sacl
	}
	if SE_DACL_PRESENT&SDDLControl(descriptor.Control) == SE_DACL_PRESENT {
		// parse DACL
		dacl, err := NewAcl(descBytes[descriptor.OffsetDacl:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse dacl bytes: %s", err)
		}
		descriptor.Dacl = dacl
	}
	return descriptor, nil
}

func (d *NtSecurityDescriptor) String() string {
	if d == nil {
		return ""
	}
	buf := strings.Builder{}
	buf.WriteString("NtSecurityDescriptor:")
	buf.WriteString("\n")
	buf.WriteString("  Revision: ")
	buf.WriteString(fmt.Sprintf("0x%02x", d.Revision))
	buf.WriteString("\n")
	buf.WriteString("  Sbz1: ")
	buf.WriteString(fmt.Sprintf("0x%02x", d.Sbz1))
	buf.WriteString("\n")
	buf.WriteString("  Control: ")
	buf.WriteString(SDDLControlToFullStr(SDDLControl(d.Control)))
	buf.WriteString("\n")
	buf.WriteString("  OwnerSid: ")
	buf.WriteString(d.OwnerSid.Alias())
	buf.WriteString("\n")
	buf.WriteString("  GroupSid: ")
	buf.WriteString(d.GroupSid.Alias())
	buf.WriteString("\n")
	buf.WriteString("  Sacl: ")
	buf.WriteString("\n")
	buf.WriteString(d.Sacl.String())
	buf.WriteString("  Dacl: ")
	buf.WriteString("\n")
	buf.WriteString(d.Dacl.String())
	return buf.String()
}

func (d *NtSecurityDescriptor) NtString() string {
	if d == nil {
		return ""
	}
	buf := strings.Builder{}
	ownerSid := d.OwnerSid.String()
	if len(ownerSid) != 0 {
		buf.WriteString("O:")
		buf.WriteString(ownerSid)
	}
	groupSid := d.GroupSid.String()
	if len(groupSid) != 0 {
		buf.WriteString("G:")
		buf.WriteString(groupSid)
	}
	if d.Dacl != nil {
		buf.WriteString("D:")
		buf.WriteString(d.Dacl.NtString())
	}
	if d.Sacl != nil {
		buf.WriteString("S:")
		buf.WriteString(d.Sacl.NtString())
	}
	return buf.String()
}

func IsInherit(flag AccessAllowedObjectFlags) bool {
	return flag&ACE_INHERITED_OBJECT_TYPE_PRESENT == ACE_INHERITED_OBJECT_TYPE_PRESENT
}

func AllowedOrDenied(aceType AceType) string {
	allowedArr := []AceType{
		ACCESS_ALLOWED_ACE_TYPE,
		ACCESS_ALLOWED_COMPOUND_ACE_TYPE,
		ACCESS_ALLOWED_OBJECT_ACE_TYPE,
		ACCESS_ALLOWED_CALLBACK_ACE_TYPE,
		ACCESS_ALLOWED_CALLBACK_OBJECT_ACE_TYPE,
	}
	deniedArr := []AceType{
		ACCESS_DENIED_ACE_TYPE,
		ACCESS_DENIED_OBJECT_ACE_TYPE,
		ACCESS_DENIED_CALLBACK_ACE_TYPE,
		ACCESS_DENIED_CALLBACK_OBJECT_ACE_TYPE,
	}
	if utils.FindIndex[AceType](allowedArr, aceType) != -1 {
		return "allowed"
	}
	if utils.FindIndex[AceType](deniedArr, aceType) != -1 {
		return "denied"
	}
	return "unkown"
}
