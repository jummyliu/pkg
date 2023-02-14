package ldap

import (
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
	"github.com/jummyliu/pkg/number"
	"github.com/jummyliu/pkg/utils"
)

// MarshalSid 把 []byte 的 objectSid 转换成可读字符串
// 	子授权机构数量（不在字符串中显示）Sub-Authority Count:1
// 	S-版本号-授权标识符-子授权机构标识符...
// 	S-{Revision:0}-{Identifier-Authority:2-8}-{Sub-Authority:8-end/4}
func MarshalSid(objectSid []byte) (string, error) {
	if len(objectSid) < 12 {
		return "", fmt.Errorf("objectSid must be at least 12 bytes long")
	}

	revision := objectSid[0]
	subAuthorityCount := int(objectSid[1])
	authorityArr := make([]uint32, subAuthorityCount)
	if len(objectSid) != 4*subAuthorityCount+2+6 {
		return "", fmt.Errorf("objectSid length is unvalid")
	}
	identifierAuthorityBytes := make([]byte, 2, 8)
	identifierAuthorityBytes = append(identifierAuthorityBytes, objectSid[2:8]...)
	// BigEndian
	identifierAuthority := binary.BigEndian.Uint64(identifierAuthorityBytes)

	begin := 8
	for i := 0; i < subAuthorityCount; i++ {
		// LittleEndian
		authorityArr[i] = binary.LittleEndian.Uint32(objectSid[begin : begin+4])
		begin += 4
	}
	sidStr := fmt.Sprintf("S-%d-%d-%s", revision, identifierAuthority, number.Join(authorityArr, "-"))
	return sidStr, nil
}

// MarshalGUID 把 []byte 的 objectGUID 转换成可读字符串
// 	生成 guid 时，前三部分字节反转
//	[0:4]-[4:6]-[6:8]-[8:10]-[10:16]
func MarshalGUID(objectGUID []byte) (string, error) {
	tmp := make([]byte, len(objectGUID))
	copy(tmp, objectGUID)
	utils.Reverse(tmp[0:4])
	utils.Reverse(tmp[4:6])
	utils.Reverse(tmp[6:8])
	guid, err := uuid.FromBytes(tmp)
	if err != nil {
		return "", fmt.Errorf("objectGUID is unvalid")
	}
	return guid.String(), nil
}
