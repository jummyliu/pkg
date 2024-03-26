package ldap

import (
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
)

// AceHeader Ace 头部信息，定义 Ace 类型和 Ace flags，4个字节
//
//	 -----------------------------------------------
//	| 0x00-0x07 | 0x08-0x0F | 0x10-0x17 | 0x18-0x1F |
//	| AceType   | AceFlags  |        AceSize        |
//	 -----------------------------------------------
type AceHeader struct {
	AceType  byte   `json:"aceType"`
	AceFlags byte   `json:"aceFlags"`
	AceSize  uint16 `json:"-"` // LittleEndian
}

// AceType Ace类型
type AceType byte

const (
	ACCESS_ALLOWED_ACE_TYPE                 AceType = 0x00
	ACCESS_DENIED_ACE_TYPE                  AceType = 0x01
	SYSTEM_AUDIT_ACE_TYPE                   AceType = 0x02
	SYSTEM_ALARM_ACE_TYPE                   AceType = 0x03
	ACCESS_ALLOWED_COMPOUND_ACE_TYPE        AceType = 0x04
	ACCESS_ALLOWED_OBJECT_ACE_TYPE          AceType = 0x05
	ACCESS_DENIED_OBJECT_ACE_TYPE           AceType = 0x06
	SYSTEM_AUDIT_OBJECT_ACE_TYPE            AceType = 0x07
	SYSTEM_ALARM_OBJECT_ACE_TYPE            AceType = 0x08
	ACCESS_ALLOWED_CALLBACK_ACE_TYPE        AceType = 0x09
	ACCESS_DENIED_CALLBACK_ACE_TYPE         AceType = 0x0a
	ACCESS_ALLOWED_CALLBACK_OBJECT_ACE_TYPE AceType = 0x0b
	ACCESS_DENIED_CALLBACK_OBJECT_ACE_TYPE  AceType = 0x0c
	SYSTEM_AUDIT_CALLBACK_ACE_TYPE          AceType = 0x0d
	SYSTEM_ALARM_CALLBACK_ACE_TYPE          AceType = 0x0e
	SYSTEM_AUDIT_CALLBACK_OBJECT_ACE_TYPE   AceType = 0x0f
	SYSTEM_ALARM_CALLBACK_OBJECT_ACE_TYPE   AceType = 0x10
	SYSTEM_MANDATORY_LABEL_ACE_TYPE         AceType = 0x11
	SYSTEM_RESOURCE_ATTRIBUTE_ACE_TYPE      AceType = 0x12
	SYSTEM_SCOPED_POLICY_ID_ACE_TYPE        AceType = 0x13
)

var AceTypeMap = map[AceType]string{
	ACCESS_ALLOWED_ACE_TYPE:                 "A",
	ACCESS_DENIED_ACE_TYPE:                  "D",
	SYSTEM_AUDIT_ACE_TYPE:                   "AU",
	SYSTEM_ALARM_ACE_TYPE:                   "AL",
	ACCESS_ALLOWED_COMPOUND_ACE_TYPE:        "", // 保留
	ACCESS_ALLOWED_OBJECT_ACE_TYPE:          "OA",
	ACCESS_DENIED_OBJECT_ACE_TYPE:           "OD",
	SYSTEM_AUDIT_OBJECT_ACE_TYPE:            "OU",
	SYSTEM_ALARM_OBJECT_ACE_TYPE:            "OL",
	ACCESS_ALLOWED_CALLBACK_ACE_TYPE:        "XA",
	ACCESS_DENIED_CALLBACK_ACE_TYPE:         "XD",
	ACCESS_ALLOWED_CALLBACK_OBJECT_ACE_TYPE: "ZA",
	ACCESS_DENIED_CALLBACK_OBJECT_ACE_TYPE:  "OD", //
	SYSTEM_AUDIT_CALLBACK_ACE_TYPE:          "XU",
	SYSTEM_ALARM_CALLBACK_ACE_TYPE:          "",   // 保留
	SYSTEM_AUDIT_CALLBACK_OBJECT_ACE_TYPE:   "OU", //
	SYSTEM_ALARM_CALLBACK_OBJECT_ACE_TYPE:   "",   // 保留
	SYSTEM_MANDATORY_LABEL_ACE_TYPE:         "ML",
	SYSTEM_RESOURCE_ATTRIBUTE_ACE_TYPE:      "RA",
	SYSTEM_SCOPED_POLICY_ID_ACE_TYPE:        "SP",
}

var AceTypeFullMap = map[AceType]string{
	ACCESS_ALLOWED_ACE_TYPE:                 "Access Allowed",
	ACCESS_DENIED_ACE_TYPE:                  "Access Denied",
	SYSTEM_AUDIT_ACE_TYPE:                   "System Audit",
	SYSTEM_ALARM_ACE_TYPE:                   "System Alarm",
	ACCESS_ALLOWED_COMPOUND_ACE_TYPE:        "Access Allowed Compound",
	ACCESS_ALLOWED_OBJECT_ACE_TYPE:          "Access Allowed Object",
	ACCESS_DENIED_OBJECT_ACE_TYPE:           "Access Denied Object",
	SYSTEM_AUDIT_OBJECT_ACE_TYPE:            "System Audit Object",
	SYSTEM_ALARM_OBJECT_ACE_TYPE:            "System Alarm Object",
	ACCESS_ALLOWED_CALLBACK_ACE_TYPE:        "Access Allowed Callback",
	ACCESS_DENIED_CALLBACK_ACE_TYPE:         "Access Denied Callback",
	ACCESS_ALLOWED_CALLBACK_OBJECT_ACE_TYPE: "Access Allowed Callback Object",
	ACCESS_DENIED_CALLBACK_OBJECT_ACE_TYPE:  "Access Denied Callback Object",
	SYSTEM_AUDIT_CALLBACK_ACE_TYPE:          "System Audit Callback",
	SYSTEM_ALARM_CALLBACK_ACE_TYPE:          "System Alarm Callback",
	SYSTEM_AUDIT_CALLBACK_OBJECT_ACE_TYPE:   "System Audit Callback Object",
	SYSTEM_ALARM_CALLBACK_OBJECT_ACE_TYPE:   "Ststem Alarm Callback Object",
	SYSTEM_MANDATORY_LABEL_ACE_TYPE:         "System Mandatory Label",
	SYSTEM_RESOURCE_ATTRIBUTE_ACE_TYPE:      "System Resource Attribute",
	SYSTEM_SCOPED_POLICY_ID_ACE_TYPE:        "System Scoped Polidy ID",
}

func AceTypeToStr(aceType AceType) string {
	if val, ok := AceTypeMap[aceType]; ok && val != "" {
		return val
	}
	return fmt.Sprintf("0x%02x", aceType)
}
func AceTypeToFullStr(aceType AceType) string {
	if val, ok := AceTypeFullMap[aceType]; ok && val != "" {
		return fmt.Sprintf("%s(0x%02x)", val, aceType)
	}
	return fmt.Sprintf("0x%02x", aceType)
}

// AceFlags ACE Flags
type AceFlags byte

const (
	OBJECT_INHERIT_ACE         AceFlags = 0x01
	CONTAINER_INHERIT_ACE      AceFlags = 0x02
	NO_PROPAGATE_INHERIT_ACE   AceFlags = 0x04
	INHERIT_ONLY_ACE           AceFlags = 0x08
	INHERITED_ACE              AceFlags = 0x10
	SUCCESSFUL_ACCESS_ACE_FLAG AceFlags = 0x40
	FAILED_ACCESS_ACE_FLAG     AceFlags = 0x80
)

var AceFlagsMap = map[AceFlags]string{
	OBJECT_INHERIT_ACE:         "OI",
	CONTAINER_INHERIT_ACE:      "CI",
	NO_PROPAGATE_INHERIT_ACE:   "NP",
	INHERIT_ONLY_ACE:           "IO",
	INHERITED_ACE:              "ID",
	SUCCESSFUL_ACCESS_ACE_FLAG: "SA",
	FAILED_ACCESS_ACE_FLAG:     "FA",
}

var AceFlagsFullMap = map[AceFlags]string{
	OBJECT_INHERIT_ACE:         "OBJECT INHERIT",
	CONTAINER_INHERIT_ACE:      "CONTAINER INHERIT",
	NO_PROPAGATE_INHERIT_ACE:   "NO PROPAGATE INHERIT",
	INHERIT_ONLY_ACE:           "INHERIT ONLY",
	INHERITED_ACE:              "INHERITED",
	SUCCESSFUL_ACCESS_ACE_FLAG: "SUCCESSFUL ACCESS ACE FLAG",
	FAILED_ACCESS_ACE_FLAG:     "FAILED ACCESS ACE FLAG",
}

func AceFlagsToStr(flags AceFlags) string {
	var tmpArr = []string{}
	for key, val := range AceFlagsMap {
		if key&flags == key {
			tmpArr = append(tmpArr, val)
		}
	}
	return strings.Join(tmpArr, "")
}

func AceFlagsToFullStr(flags AceFlags) string {
	var tmpArr = []string{}
	for key, val := range AceFlagsFullMap {
		if key&flags == key {
			tmpArr = append(tmpArr, fmt.Sprintf("%s(0x%02x)", val, key))
		}
	}
	return strings.Join(tmpArr, ",")
}

func GetInheritanceFlags(flags AceFlags) string {
	result := "None"
	arr := []AceFlags{CONTAINER_INHERIT_ACE, OBJECT_INHERIT_ACE}
	for _, key := range arr {
		if key&flags == key {
			result = fmt.Sprintf("%s(0x%02x)", AceFlagsFullMap[key], key)
		}
	}
	return result
}

func GetPropagationFlags(flags AceFlags) string {
	result := "None"
	arr := []AceFlags{NO_PROPAGATE_INHERIT_ACE, INHERIT_ONLY_ACE}
	for _, key := range arr {
		if key&flags == key {
			result = fmt.Sprintf("%s(0x%02x)", AceFlagsFullMap[key], key)
		}
	}
	return result
}

type AceMask uint32

const (
	// generic rights 是抽象的权限，会根据不同的对象类型，映射不同的权限
	ADS_RIGHT_GENERIC_READ    AceMask = 0x80000000 // 读
	ADS_RIGHT_GENERIC_WRITE   AceMask = 0x40000000 // 写
	ADS_RIGHT_GENERIC_EXECUTE AceMask = 0x20000000 // 列出容器内容的权限
	ADS_RIGHT_GENERIC_ALL     AceMask = 0x10000000 // 所有权限

	GENERIC_READ    AceMask = 131220 // 实际 GENERIC_READ 的掩码
	GENERIC_WRITE   AceMask = 131112 // 实际 GENERIC_WRITE 的掩码
	GENERIC_EXECUTE AceMask = 131076 // 实际 GENERIC_EXECUTE 的掩码
	GENERIC_ALL     AceMask = 983551 // 实际 GENERIC_ALL 的掩码

	ADS_RIGHT_MAXIMUM_ALLOWED        AceMask = 0x02000000
	ADS_RIGHT_ACCESS_SYSTEM_SECURITY AceMask = 0x01000000 // 读写SACL权限
	ADS_RIGHT_SYNCHRONIZE            AceMask = 0x00100000 // 同步的权限

	// std rights
	ADS_RIGHT_WRITE_OWNER  AceMask = 0x00080000 // 所有者的权限
	ADS_RIGHT_WRITE_DAC    AceMask = 0x00040000 // 修改DACL权限
	ADS_RIGHT_READ_CONTROL AceMask = 0x00020000 // 读ntSecurityDescriptor权限（不含SACL）
	ADS_RIGHT_DELETE       AceMask = 0x00010000 // 删除权限

	// ds right
	ADS_RIGHT_DS_CREATE_CHILD   AceMask = 0x00000001 // 新建子对象的权限
	ADS_RIGHT_DS_DELETE_CHILD   AceMask = 0x00000002 // 删除子对象的权限
	ADS_RIGHT_ACTRL_DS_LIST     AceMask = 0x00000004 // 列出自对象的权限
	ADS_RIGHT_DS_SELF           AceMask = 0x00000008
	ADS_RIGHT_DS_READ_PROP      AceMask = 0x00000010 // 读属性
	ADS_RIGHT_DS_WRITE_PROP     AceMask = 0x00000020 // 写属性
	ADS_RIGHT_DS_DELETE_TREE    AceMask = 0x00000040 // 删除子对象
	ADS_RIGHT_DS_LIST_OBJECT    AceMask = 0x00000080 // 列出对象权限
	ADS_RIGHT_DS_CONTROL_ACCESS AceMask = 0x00000100
)

var AceMasksMap = map[AceMask]string{
	ADS_RIGHT_GENERIC_READ:           "GR",
	ADS_RIGHT_GENERIC_WRITE:          "GW",
	ADS_RIGHT_GENERIC_EXECUTE:        "GX",
	ADS_RIGHT_GENERIC_ALL:            "GA",
	ADS_RIGHT_MAXIMUM_ALLOWED:        "MA",
	ADS_RIGHT_ACCESS_SYSTEM_SECURITY: "AS",
	ADS_RIGHT_SYNCHRONIZE:            "SY",
	ADS_RIGHT_WRITE_OWNER:            "WO",
	ADS_RIGHT_WRITE_DAC:              "WD",
	ADS_RIGHT_READ_CONTROL:           "RC",
	ADS_RIGHT_DELETE:                 "DE",
	ADS_RIGHT_DS_CREATE_CHILD:        "CC",
	ADS_RIGHT_DS_DELETE_CHILD:        "DC",
	ADS_RIGHT_ACTRL_DS_LIST:          "LC",
	ADS_RIGHT_DS_SELF:                "SW",
	ADS_RIGHT_DS_READ_PROP:           "RP",
	ADS_RIGHT_DS_WRITE_PROP:          "WP",
	ADS_RIGHT_DS_DELETE_TREE:         "DT",
	ADS_RIGHT_DS_LIST_OBJECT:         "LO",
	ADS_RIGHT_DS_CONTROL_ACCESS:      "CR",
}
var AceMasksFullMap = map[AceMask]string{
	ADS_RIGHT_GENERIC_READ:           "GENERIC_READ(bit)",
	ADS_RIGHT_GENERIC_WRITE:          "GENERIC_WRITE(bit)",
	ADS_RIGHT_GENERIC_EXECUTE:        "GENERIC_EXECUTE(bit)",
	ADS_RIGHT_GENERIC_ALL:            "GENERIC_ALL(bit)",
	GENERIC_READ:                     "GENERIC_READ(mask)",
	GENERIC_WRITE:                    "GENERIC_WRITE(mask)",
	GENERIC_EXECUTE:                  "GENERIC_EXECUTE(mask)",
	GENERIC_ALL:                      "GENERIC_ALL(mask)",
	ADS_RIGHT_MAXIMUM_ALLOWED:        "MAXIMUM_ALLOWED",
	ADS_RIGHT_ACCESS_SYSTEM_SECURITY: "ACCESS_SYSTEM_SECURITY",
	ADS_RIGHT_SYNCHRONIZE:            "SYNCHRONIZE",
	ADS_RIGHT_WRITE_OWNER:            "WRITE_OWNER",
	ADS_RIGHT_WRITE_DAC:              "WRITE_DAC",
	ADS_RIGHT_READ_CONTROL:           "READ_CONTROL",
	ADS_RIGHT_DELETE:                 "DELETE",
	ADS_RIGHT_DS_CREATE_CHILD:        "DS_CREATE_CHILD",
	ADS_RIGHT_DS_DELETE_CHILD:        "DS_DELETE_CHILD",
	ADS_RIGHT_ACTRL_DS_LIST:          "ACTRL_DS_LIST",
	ADS_RIGHT_DS_SELF:                "DS_SELF",
	ADS_RIGHT_DS_READ_PROP:           "DS_READ_PROP",
	ADS_RIGHT_DS_WRITE_PROP:          "DS_WRITE_PROP",
	ADS_RIGHT_DS_DELETE_TREE:         "DS_DELETE_TREE",
	ADS_RIGHT_DS_LIST_OBJECT:         "DS_LIST_OBJECT",
	ADS_RIGHT_DS_CONTROL_ACCESS:      "DS_CONTROL_ACCESS",
}

func AceMaskToStr(mask AceMask) string {
	var tmpArr = []string{}
	for key, val := range AceMasksMap {
		if key&mask == key {
			tmpArr = append(tmpArr, val)
		}
	}
	return strings.Join(tmpArr, "")
}
func AceMaskToFullStr(mask AceMask) string {
	var tmpArr = []string{}
	for key, val := range AceMasksFullMap {
		if key&mask == key {
			tmpArr = append(tmpArr, fmt.Sprintf("%s(0x%08x)", val, key))
		}
	}
	return strings.Join(tmpArr, ",")
}

// ACCESS_ALLOWED_OBJECT_ACE, ACCESS_DENIED_OBJECT_ACE 的 flags
type AccessAllowedObjectFlags uint32

const (
	Null                              AccessAllowedObjectFlags = 0x00000000
	ACE_OBJECT_TYPE_PRESENT           AccessAllowedObjectFlags = 0x00000001
	ACE_INHERITED_OBJECT_TYPE_PRESENT AccessAllowedObjectFlags = 0x00000002
)

var AccessAllowedObjectFlagsMap = map[AccessAllowedObjectFlags]string{
	// Null:                              "Null",
	ACE_OBJECT_TYPE_PRESENT:           "ACE_OBJECT_TYPE_PRESENT",
	ACE_INHERITED_OBJECT_TYPE_PRESENT: "ACE_INHERITED_OBJECT_TYPE_PRESENT",
}

func AccessAllowedObjectFlagsToStr(flags AccessAllowedObjectFlags) string {
	var tmpArr = []string{}
	for key, val := range AccessAllowedObjectFlagsMap {
		if key&flags == key {
			tmpArr = append(tmpArr, val)
		}
	}
	return strings.Join(tmpArr, " ")
}
func AccessAllowedObjectFlagsToFullStr(flags AccessAllowedObjectFlags) string {
	var tmpArr = []string{}
	for key, val := range AccessAllowedObjectFlagsMap {
		if key&flags == key {
			tmpArr = append(tmpArr, fmt.Sprintf("%s(0x%08x)", val, key))
		}
	}
	return strings.Join(tmpArr, ",")
}

type Ace interface {
	Size() int // 返回 Ace 大小，即下一个的偏移量
	NtString() string
	fmt.Stringer
}

// NewAce 创建 Ace
func NewAce(aceBytes []byte) (ace Ace, err error) {
	if len(aceBytes) < 4 {
		return nil, fmt.Errorf("ace must be at least 4 bytes long")
	}
	aceType := aceBytes[0]
	// aceHeader := &AceHeader{
	// 	AceType:  aceBytes[0],
	// 	AceFlags: aceBytes[1],
	// 	AceSize:  binary.LittleEndian.Uint16(aceBytes[2:4]),
	// }
	aceConstructorsMu.RLock()
	constructor, ok := aceConstructors[AceType(aceType)]
	aceConstructorsMu.RUnlock()
	if !ok {
		// 没有注册，使用默认解析器
		return NewDefaultAce(aceBytes)
		// return nil, fmt.Errorf("ace: unknown constructor for aceType %d", aceType)
	}
	return constructor(aceBytes)
}

// AceConstructor Ace 构造函数签名
type AceConstructor func(aceBytes []byte) (Ace, error)

var (
	aceConstructorsMu sync.RWMutex
	aceConstructors   = make(map[AceType]AceConstructor)
)

// RegisterAce 注册 Ace 解析构造函数
func RegisterAce(aceType AceType, constructor AceConstructor) {
	aceConstructorsMu.Lock()
	defer aceConstructorsMu.Unlock()
	if aceConstructors == nil {
		panic("ace: Register aceConstructors is nil")
	}
	if _, ok := aceConstructors[aceType]; ok {
		panic(fmt.Sprintf("ace: Register called twice for aceConstructors %d", aceType))
	}
	aceConstructors[aceType] = constructor
}

// DefaultAce 默认 Ace 结构体，只解析 AceHeader
type DefaultAce struct {
	*AceHeader
	Data []byte `json:"-"`
}

func NewDefaultAce(aceBytes []byte) (ace Ace, err error) {
	if len(aceBytes) < 4 {
		return nil, fmt.Errorf("ace must be at least 4 bytes long")
	}
	aceHeader := &AceHeader{
		AceType:  aceBytes[0],
		AceFlags: aceBytes[1],
		AceSize:  binary.LittleEndian.Uint16(aceBytes[2:4]),
	}
	if len(aceBytes) < int(aceHeader.AceSize) {
		return nil, fmt.Errorf("ace length is unvalid")
	}
	a := &DefaultAce{
		AceHeader: aceHeader,
		Data:      aceBytes[4:],
	}
	return a, nil
}

func (a *DefaultAce) Size() int {
	if a == nil {
		return 0
	}
	return int(a.AceHeader.AceSize)
}

func (a *DefaultAce) String() string {
	if a == nil {
		return ""
	}
	buf := strings.Builder{}
	buf.WriteString("        AceType: ")
	buf.WriteString(AceTypeToFullStr(AceType(a.AceHeader.AceType)))
	buf.WriteString("\n")
	buf.WriteString("          ")
	buf.WriteString(AllowedOrDenied(AceType(a.AceHeader.AceType)))
	buf.WriteString("\n")
	buf.WriteString("        HeaderAceFlags: ")
	buf.WriteString(fmt.Sprintf("0x%02x", a.AceHeader.AceFlags))
	buf.WriteString("\n")
	buf.WriteString("          ")
	buf.WriteString(AceFlagsToFullStr(AceFlags(a.AceHeader.AceFlags)))
	buf.WriteString("\n")
	buf.WriteString("        AceSize: ")
	buf.WriteString(fmt.Sprintf("%d", a.AceHeader.AceSize))
	buf.WriteString("\n")
	buf.WriteString("        Ace Data: unresolved")
	buf.WriteString("\n")
	return buf.String()
}

func (a *DefaultAce) NtString() string {
	return ""
}
