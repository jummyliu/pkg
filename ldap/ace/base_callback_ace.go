package ace

import "github.com/jummyliu/pkg/ldap"

// BaseCallbackAce
//
//	ACCESS_ALLOWED_CALLBACK_ACE, ACCESS_DENIED_CALLBACK_ACE, SYSTEM_AUDIT_CALLBACK_ACE
//
//	  -----------------------------------------------
//	| 0x00-0x07 | 0x08-0x0F | 0x10-0x17 | 0x18-0x1F |
//	|                    Header                     |
//	|                     Mask                      |
//	|                 Sid(variable)                 |
//	|           ApplicationData(variable)           |
//	 -----------------------------------------------
type BaseCallbackAce struct {
	Header          ldap.AceHeader // [4]byte
	Mask            [4]byte
	Sid             []byte
	ApplicationData []byte
}

func NewBaseCallbackAce(aceBytes []byte) (ace ldap.Ace, err error) {
	return nil, nil
}

func init() {
	arr := []ldap.AceType{
		ldap.ACCESS_ALLOWED_CALLBACK_ACE_TYPE,
		ldap.ACCESS_DENIED_CALLBACK_ACE_TYPE,
		ldap.SYSTEM_AUDIT_CALLBACK_ACE_TYPE,
	}
	for _, item := range arr {
		// ldap.RegisterAce(item, NewBaseCallbackAce)
		ldap.RegisterAce(item, NewBaseAce)
	}
}
