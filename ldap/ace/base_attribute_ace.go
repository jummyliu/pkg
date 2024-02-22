package ace

import "github.com/jummyliu/pkg/ldap"

// BaseAttributeAce
//
//	SYSTEM_RESOURCE_ATTRIBUTE_ACE
//
//	  -----------------------------------------------
//	| 0x00-0x07 | 0x08-0x0F | 0x10-0x17 | 0x18-0x1F |
//	|                    Header                     |
//	|                     Mask                      |
//	|                 Sid(variable)                 |
//	|           Attribute Data(variable)            |
//	 -----------------------------------------------
type BaseAttributeAce struct {
	Header        ldap.AceHeader // [4]byte
	Mask          [4]byte
	Sid           []byte
	AttributeData []byte
}

func NewBaseAttributeAce(aceBytes []byte) (ace ldap.Ace, err error) {
	return nil, nil
}

func init() {
	arr := []ldap.AceType{
		ldap.SYSTEM_RESOURCE_ATTRIBUTE_ACE_TYPE,
	}
	for _, item := range arr {
		// ldap.RegisterAce(item, NewBaseAttributeAce)
		ldap.RegisterAce(item, NewBaseAce)
	}
}
