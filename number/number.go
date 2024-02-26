package number

import (
	"fmt"
	"strconv"
	"strings"
)

type intnumber interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type floatnumber interface {
	float32 | float64
}

type number interface {
	intnumber | floatnumber
}

// ParseInt 把 any 解析成指定的 int number 类型
//
//	T 直接返回
//	数字类型，先转成字符串，再使用 strconv 转换
//	字符串，使用 strconv 转换
//	其他，返回 0
func ParseInt[T intnumber](val any) (result T) {
	str := ""
	switch v := val.(type) {
	case T:
		return v
	case string:
		str = v
	case int:
		return T(v)
	case uint:
		return T(v)
	case int8:
		return T(v)
	case uint8:
		return T(v)
	case int16:
		return T(v)
	case uint16:
		return T(v)
	case int32:
		return T(v)
	case uint32:
		return T(v)
	case int64:
		return T(v)
	case uint64:
		return T(v)
	case float32:
		return T(v)
	case float64:
		return T(v)
	// case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
	// 	str = fmt.Sprintf("%v", v)
	default:
		result = 0
	}
	if len(str) > 0 {
		tmp, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			tmp = 0
		}
		result = T(tmp)
	}
	return result
}

// ParseFloat 把 any 解析成指定的 float number 类型
//
//	T 直接返回
//	数字类型，先转成字符串，再使用 strconv 转换
//	字符串，使用 strconv 转换
//	其他，返回 0
func ParseFloat[T floatnumber](val any) (result T) {
	str := ""
	switch v := val.(type) {
	case T:
		return v
	case string:
		str = v
	case int:
		return T(v)
	case uint:
		return T(v)
	case int8:
		return T(v)
	case uint8:
		return T(v)
	case int16:
		return T(v)
	case uint16:
		return T(v)
	case int32:
		return T(v)
	case uint32:
		return T(v)
	case int64:
		return T(v)
	case uint64:
		return T(v)
	case float32:
		return T(v)
	case float64:
		return T(v)
	// case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
	// 	str = fmt.Sprintf("%v", v)
	default:
		result = 0
	}
	if len(str) > 0 {
		tmp, err := strconv.ParseFloat(str, 64)
		if err != nil {
			tmp = 0
		}
		result = T(tmp)
	}
	return result
}

// Join 使用指定分隔符 sep 连接 []number
func Join[T number](elems []T, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%v", elems[0])
	}
	// 长度应该不够
	n := len(sep)*(len(elems)-1) + len(elems)
	var b strings.Builder
	b.Grow(n)
	b.WriteString(fmt.Sprintf("%v", elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(fmt.Sprintf("%v", s))
	}
	return b.String()
}
