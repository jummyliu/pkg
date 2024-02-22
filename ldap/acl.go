package ldap

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Acl
//
//	 -----------------------------------------------
//	| 0x00-0x07 | 0x08-0x0F | 0x10-0x17 | 0x18-0x1F |
//	| Revision  | Sbz1      |        AclSize        |
//	|        AceCount       |         Sbz2          |
//	|                 AceList(variable)             |
//	 -----------------------------------------------
type Acl struct {
	AclRevision byte   `json:"-"`
	Sbz1        byte   `json:"-"`
	AclSize     uint16 `json:"-"`        // LittleEndian
	AceCount    uint16 `json:"aceCount"` // LittleEndian
	Sbz2        uint16 `json:"-"`        // LittleEndian
	AceList     []Ace  `json:"aceList"`
}

type AclRevision byte

const (
	ACL_REVISION    AclRevision = 0x02 // AceTypes 仅允许 0x00, 0x01, 0x02, 0x03, 0x11, 0x12, 0x13
	ACL_REVISION_DS AclRevision = 0x04 // AceTypes 仅允许 0x05, 0x06, 0x07, 0x08, 0x11
)

func NewAcl(aclBytes []byte) (acl *Acl, err error) {
	if len(aclBytes) < 8 {
		return nil, fmt.Errorf("acl must be at least 8 bytes long")
	}
	acl = &Acl{
		AclRevision: aclBytes[0],
		Sbz1:        aclBytes[1],
		AclSize:     binary.LittleEndian.Uint16(aclBytes[2:4]),
		AceCount:    binary.LittleEndian.Uint16(aclBytes[4:6]),
		Sbz2:        binary.LittleEndian.Uint16(aclBytes[6:8]),
	}
	if len(aclBytes) < int(acl.AclSize) {
		return nil, fmt.Errorf("acl length is unvalid")
	}
	aceList := make([]Ace, acl.AceCount)
	aceBytes := aclBytes[8:]
	for i := 0; i < int(acl.AceCount); i++ {
		ace, err := NewAce(aceBytes)
		if ace == nil || err != nil {
			return nil, fmt.Errorf("failed to parse ace bytes: %s", err)
		}
		aceList[i] = ace
		size := ace.Size()
		// 偏移
		aceBytes = aceBytes[size:]
	}
	acl.AceList = aceList
	return acl, nil
}

func (a *Acl) String() string {
	if a == nil {
		return ""
	}
	buf := strings.Builder{}
	buf.WriteString("    AclRevision: ")
	buf.WriteString(fmt.Sprintf("0x%02x", a.AclRevision))
	buf.WriteString("\n")
	buf.WriteString("    Sbz1: ")
	buf.WriteString(fmt.Sprintf("0x%02x", a.Sbz1))
	buf.WriteString("\n")
	buf.WriteString("    AclSize: ")
	buf.WriteString(fmt.Sprintf("%d", a.AclSize))
	buf.WriteString("\n")
	buf.WriteString("    AceCount: ")
	buf.WriteString(fmt.Sprintf("%d", a.AceCount))
	buf.WriteString("\n")
	buf.WriteString("    Sbz2: ")
	buf.WriteString(fmt.Sprintf("0x%02x", a.Sbz2))
	buf.WriteString("\n")
	buf.WriteString("    AceList: ")
	buf.WriteString("\n")
	for i, item := range a.AceList {
		buf.WriteString(fmt.Sprintf("      [%d]", i))
		buf.WriteString("\n")
		buf.WriteString(item.String())
	}
	return buf.String()
}

func (a *Acl) NtString() string {
	if a == nil {
		return ""
	}
	buf := strings.Builder{}
	for _, item := range a.AceList {
		buf.WriteString(item.NtString())
	}
	return buf.String()
}
