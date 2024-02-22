package ldap

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/jummyliu/pkg/number"
)

// Sid
//
//	struct:
//	 -------------------------------------------------------------------------------
//	| 0x00-0x07         | 0x08-0x0F         | 0x10-0x17         | 0x18-0x1F         |
//	| Revision(0x01)    | SubAuthorityCount |      IdentifierAuthority(6 bytes)     |
//	|                               SubAuthority(variable)                          |
//	 -------------------------------------------------------------------------------
//
//	string:
//		子授权机构数量（不在字符串中显示）SubAuthorityCount
//		S-版本号-授权标识符-子授权机构标识符...
//		S-{Revision}-{IdentifierAUthority}-{SubAuthority...}
type Sid struct {
	Revision            byte
	SubAuthorityCount   byte     // maximum number is 15
	IdentifierAuthority [6]byte  // BigEndian
	SubAuthority        []uint32 // size is SubAuthorityCount
}

// NewSid 把 []byte 的 Sid 转换成 Sid 结构体指针
//
//	如果转换失败，则返回错误，并且 Sid 为 nil
func NewSid(sidBytes []byte) (*Sid, error) {
	if len(sidBytes) < 12 {
		return nil, fmt.Errorf("sid must be at least 12 bytes long")
	}
	revision := sidBytes[0]
	subAuthorityCount := sidBytes[1]
	if subAuthorityCount > 15 || len(sidBytes) < 4*int(subAuthorityCount)+2+6 {
		return nil, fmt.Errorf("sid length is unvalid")
	}
	var identifierAuthority [6]byte
	copy(identifierAuthority[:], sidBytes[2:8])
	subAuthority := make([]uint32, subAuthorityCount)
	err := binary.Read(bytes.NewReader(sidBytes[8:4*subAuthorityCount+8]), binary.LittleEndian, &subAuthority)
	if err != nil {
		return nil, fmt.Errorf("fialed to parse the SubAuthority property of the sid: %s", err)
	}

	return &Sid{
		Revision:            revision,
		SubAuthorityCount:   subAuthorityCount,
		IdentifierAuthority: identifierAuthority,
		SubAuthority:        subAuthority,
	}, nil
}

// String implement fmt.Stringer interface
func (s *Sid) String() string {
	if s == nil {
		return ""
	}
	identifierAuthorityBytes := make([]byte, 2, 8)
	identifierAuthorityBytes = append(identifierAuthorityBytes, s.IdentifierAuthority[:]...)
	identifierAuthority := binary.BigEndian.Uint64(identifierAuthorityBytes)
	return fmt.Sprintf("S-%d-%d-%s", s.Revision, identifierAuthority, number.Join(s.SubAuthority, "-"))
}

// Alias return well known sid
func (s *Sid) Alias() string {
	sid := s.String()
	if val, ok := WellKnownSids[sid]; ok {
		return fmt.Sprintf("%s(%s)", val, sid)
	}
	return sid
}

func (s Sid) MarshalJSON() ([]byte, error) {
	tmp := fmt.Sprintf("\"%s\"", s.String())
	return []byte(tmp), nil
}

var WellKnownSids = map[string]string{
	"S-1-0":        "BUILTIN\\Null Authority",
	"S-1-0-0":      "BUILTIN\\Nobody",
	"S-1-1":        "BUILTIN\\World Authority",
	"S-1-1-0":      "BUILTIN\\Everyone",
	"S-1-2":        "BUILTIN\\Local Authority",
	"S-1-2-0":      "BUILTIN\\Local",
	"S-1-2-1":      "BUILTIN\\Console Logon",
	"S-1-3":        "BUILTIN\\Creator Authority",
	"S-1-3-0":      "BUILTIN\\Creator Owner",
	"S-1-3-1":      "BUILTIN\\Creator Group",
	"S-1-3-2":      "BUILTIN\\Creator Owner Server",
	"S-1-3-3":      "BUILTIN\\Creator Group Server",
	"S-1-3-4":      "BUILTIN\\Owner Rights",
	"S-1-5-80-0":   "BUILTIN\\All Services",
	"S-1-4":        "BUILTIN\\Non-unique Authority",
	"S-1-5":        "BUILTIN\\NT Authority",
	"S-1-5-1":      "BUILTIN\\Dialup",
	"S-1-5-2":      "BUILTIN\\Network",
	"S-1-5-3":      "BUILTIN\\Batch",
	"S-1-5-4":      "BUILTIN\\Interactive",
	"S-1-5-6":      "BUILTIN\\Service",
	"S-1-5-7":      "BUILTIN\\Anonymous",
	"S-1-5-8":      "BUILTIN\\Proxy",
	"S-1-5-9":      "BUILTIN\\Enterprise Domain Controllers",
	"S-1-5-10":     "BUILTIN\\Principal Self",
	"S-1-5-11":     "BUILTIN\\Authenticated Users",
	"S-1-5-12":     "BUILTIN\\Restricted Code",
	"S-1-5-13":     "BUILTIN\\Terminal Server Users",
	"S-1-5-14":     "BUILTIN\\Remote Interactive Logon",
	"S-1-5-15":     "BUILTIN\\This Organization",
	"S-1-5-17":     "BUILTIN\\This Organization",
	"S-1-5-18":     "BUILTIN\\Local System",
	"S-1-5-19":     "BUILTIN\\NT Authority",
	"S-1-5-20":     "BUILTIN\\NT Authority",
	"S-1-5-80":     "BUILTIN\\NT Service",
	"S-1-5-83-0":   "NT VIRTUAL MACHINE\\Virtual Machines",
	"S-1-16-0":     "BUILTIN\\Untrusted Mandatory Level",
	"S-1-5-32-544": "BUILTIN\\Administrators",
	"S-1-5-32-545": "BUILTIN\\Users",
	"S-1-5-32-546": "BUILTIN\\Guests",
	"S-1-5-32-547": "BUILTIN\\Power Users",
	"S-1-5-32-548": "BUILTIN\\Account Operators",
	"S-1-5-32-549": "BUILTIN\\Server Operators",
	"S-1-5-32-550": "BUILTIN\\Print Operators",
	"S-1-5-32-551": "BUILTIN\\Backup Operators",
	"S-1-5-32-552": "BUILTIN\\Replicators",
	"S-1-5-64-10":  "BUILTIN\\NTLM Authentication",
	"S-1-5-64-14":  "BUILTIN\\SChannel Authentication",
	"S-1-5-64-21":  "BUILTIN\\Digest Authentication",
	"S-1-16-4096":  "BUILTIN\\Low Mandatory Level",
	"S-1-16-8192":  "BUILTIN\\Medium Mandatory Level",
	"S-1-16-8448":  "BUILTIN\\Medium Plus Mandatory Level",
	"S-1-16-12288": "BUILTIN\\High Mandatory Level",
	"S-1-16-16384": "BUILTIN\\System Mandatory Level",
	"S-1-16-20480": "BUILTIN\\Protected Process Mandatory Level",
	"S-1-16-28672": "BUILTIN\\Secure Process Mandatory Level",
	"S-1-5-32-554": "BUILTIN\\Pre-Windows 2000 Compatible Access",
	"S-1-5-32-555": "BUILTIN\\Remote Desktop Users",
	"S-1-5-32-556": "BUILTIN\\Network Configuration Operators",
	"S-1-5-32-557": "BUILTIN\\Incoming Forest Trust Builders",
	"S-1-5-32-558": "BUILTIN\\Performance Monitor Users",
	"S-1-5-32-559": "BUILTIN\\Performance Log Users",
	"S-1-5-32-560": "BUILTIN\\Windows Authorization Access Group",
	"S-1-5-32-561": "BUILTIN\\Terminal Server License Servers",
	"S-1-5-32-562": "BUILTIN\\Distributed COM Users",
	"S-1-5-32-569": "BUILTIN\\Cryptographic Operators",
	"S-1-5-32-573": "BUILTIN\\Event Log Readers",
	"S-1-5-32-574": "BUILTIN\\Certificate Service DCOM Access",
	"S-1-5-32-575": "BUILTIN\\RDS Remote Access Servers",
	"S-1-5-32-576": "BUILTIN\\RDS Endpoint Servers",
	"S-1-5-32-577": "BUILTIN\\RDS Management Servers",
	"S-1-5-32-578": "BUILTIN\\Hyper-V Administrators",
	"S-1-5-32-579": "BUILTIN\\Access Control Assistance Operators",
	"S-1-5-32-580": "BUILTIN\\Remote Management Users",
}
