package ace

type LiteralToken byte

const (
	INVALID_TOKEN  LiteralToken = 0x00
	SIGNED_INT8    LiteralToken = 0x01
	SIGNED_INT16   LiteralToken = 0x02
	SIGNED_INT32   LiteralToken = 0x03
	SIGNED_INT64   LiteralToken = 0x04
	UNICODE_STRING LiteralToken = 0x10
	OCTET_STRING   LiteralToken = 0x18
	COMPOSITE      LiteralToken = 0x50
	SID            LiteralToken = 0x51
)

const (
	OCTAL       LiteralToken = 0x01
	DECIMAL     LiteralToken = 0x02
	HEXADECIMAL LiteralToken = 0x03
)

const (
	PLUS    LiteralToken = 0x01
	MINUS   LiteralToken = 0x02
	NO_SIGN LiteralToken = 0x03
)

const (
	// Unary Relational Operators
	MEMBER_OF                LiteralToken = 0x89
	DEVICE_MEMBER_OF         LiteralToken = 0x8a
	MEMBER_OF_ANY            LiteralToken = 0x8b
	DEVICE_MEMBER_OF_ANY     LiteralToken = 0x8c
	NOT_MEMBER_OF            LiteralToken = 0x90
	NOT_DEVICE_MEMBER_OF     LiteralToken = 0x91
	NOT_MEMBER_OF_ANY        LiteralToken = 0x92
	NOT_DEIVCE_MEMBER_OF_ANY LiteralToken = 0x93

	// Binary Relational Operators
	EQUAL        LiteralToken = 0x80
	NO_EQUAL     LiteralToken = 0x81
	LT           LiteralToken = 0x82
	LTE          LiteralToken = 0x83
	GT           LiteralToken = 0x84
	GTE          LiteralToken = 0x85
	CONTAINS     LiteralToken = 0x86
	ANY_OF       LiteralToken = 0x88
	NOT_CONTAINS LiteralToken = 0x8e
	NOT_ANY_OF   LiteralToken = 0x8f

	// Logical Operator Tokens
	EXISTS      LiteralToken = 0x87
	NOT_EXISTS  LiteralToken = 0x8d
	LOGICAL_NOT LiteralToken = 0xa2
	// Binary Logical Operators
	LOGICAL_AND LiteralToken = 0xa0 // &&
	LOGICAL_OR  LiteralToken = 0xa1 // ||
)

// Attribute Tokens
const (
	LOCAL_ATTRIBUTE    LiteralToken = 0xf8 // ||
	USER_ATTRIBUTE     LiteralToken = 0xf9 // ||
	RESOURCE_ATTRIBUTE LiteralToken = 0xfa // ||
	DEVICE_ATTRIBUTE   LiteralToken = 0xfb // ||
)

// ApplicationData 条件ACE
//
//	(Title=="VP")
//	 ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//	|                           | Local Attribute "Title"                                             | Literal value "VP"                                            | Operator "=="  |          |
//	| Conditional-ace signature | Attribute token                                                     | String literal token                                          | "==" token     | Padding  |
//	| Signature bytes(4 bytes)  | Attribute byte-code | Length(DWORD) | Unicode characters            | Unicode string byte-code | Length(DWORD) | Unicode characters | "==" byte-code |          |
//	| 61 72 74 78               | f8                  | a 0 0 0       | 54 00 69 00 74 00 6c 00 65 00 | 10                       | 4 0 0 0       | 56 00 50 00        | 80             | 00 00 00 |
//	 ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
type ApplicationData struct {
}
