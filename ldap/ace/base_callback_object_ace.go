package ace

import "github.com/jummyliu/pkg/ldap"

// BaseCallbackObjectAce
//
//	ACCESS_ALLOWED_CALLBACK_OBJECT_ACE, ACCESS_DENIED_CALLBACK_OBJECT_ACE, SYSTEM_AUDIT_OBJECT_ACE, SYSTEM_AUDIT_CALLBACK_OBJECT_ACE
//
//	  -----------------------------------------------
//	| 0x00-0x07 | 0x08-0x0F | 0x10-0x17 | 0x18-0x1F |
//	|                    Header                     |
//	|                     Mask                      |
//	|                    Flags                      |
//	|              ObjectType(16 bytes)             |
//	|         InheritedObjectType(16 bytes)         |
//	|                 Sid(variable)                 |
//	|           ApplicationData(variable)           |
//	 -----------------------------------------------
type BaseCallbackObjectAce struct {
	Header              ldap.AceHeader // [4]byte
	Mask                [4]byte
	Flags               [4]byte
	ObjectType          [16]byte
	InheritedObjectType [16]byte
	Sid                 []byte
	ApplicationData     []byte
}

func NewBaseCallbackObjectAce(aceBytes []byte) (ace ldap.Ace, err error) {
	return nil, nil
}

func init() {
	arr := []ldap.AceType{
		ldap.ACCESS_ALLOWED_CALLBACK_OBJECT_ACE_TYPE,
		ldap.ACCESS_DENIED_CALLBACK_OBJECT_ACE_TYPE,
		ldap.SYSTEM_AUDIT_OBJECT_ACE_TYPE,
		ldap.SYSTEM_AUDIT_CALLBACK_OBJECT_ACE_TYPE,
	}
	for _, item := range arr {
		// ldap.RegisterAce(item, NewBaseCallbackObjectAce)
		ldap.RegisterAce(item, NewBaseObjectAce)
	}
}
