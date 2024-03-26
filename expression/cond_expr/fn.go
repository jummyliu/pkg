package cond_expr

import (
	"regexp"
	"strings"

	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/jummyliu/pkg/expression/token"
	"github.com/jummyliu/pkg/number"
	"github.com/jummyliu/pkg/utils"
)

// typeCheck 类型检查
func typeCheck[T comparable](m map[string]any, key string, value any) (val T, ok bool) {
	v, err := mapstr.M(m).GetValue(key)
	if err != nil {
		return val, false
	}
	vNew, ok := v.(T)
	if !ok {
		return val, false
	}
	if _, ok := value.(T); !ok {
		return val, false
	}
	return vNew, true
}

type ConditionFn func(m map[string]any, key string, value any) bool

var DefaultFnMap = map[string]map[token.Token]ConditionFn{
	"==": {
		token.NUM:    equal[float64],
		token.BOOL:   equal[bool],
		token.STRING: equalStr,
	},
	"!=": {
		token.NUM:    unEqual[float64],
		token.BOOL:   unEqual[bool],
		token.STRING: unEqualStr,
	},
	">=": {
		token.NUM:    gte[float64],
		token.STRING: gte[string],
	},
	"<=": {
		token.NUM:    lte[float64],
		token.STRING: lte[string],
	},
	">": {
		token.NUM:    gt[float64],
		token.STRING: lt[string],
	},
	"<": {
		token.NUM:    lt[float64],
		token.STRING: lt[string],
	},
	"contains": {
		token.STRING: contains,
	},
	"unContains": {
		token.STRING: unContains,
	},
	"startsWith": {
		token.STRING: startsWith,
	},
	"unStartsWith": {
		token.STRING: unStartsWith,
	},
	"endsWith": {
		token.STRING: endsWith,
	},
	"unEndsWith": {
		token.STRING: unEndsWith,
	},
	"reg": {
		token.STRING: reg,
	},
	"in": {
		token.STRING: in,
	},
	"notIn": {
		token.STRING: notIn,
	},
	"containsBit": {
		token.NUM:    containsBit,
		token.STRING: containsBit,
	},
	"unContainsBit": {
		token.NUM:    unContainsBit,
		token.STRING: unContainsBit,
	},
	"&": {},
	"|": {},
}

// equal ==
func equal[T comparable](m map[string]any, key string, value any) bool {
	val, ok := typeCheck[T](m, key, value)
	if !ok {
		return false
	}
	return val == value
}

// equalStr ==
func equalStr(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return strings.EqualFold(val, value.(string))
}

// unEqual !=
func unEqual[T comparable](m map[string]any, key string, value any) bool {
	val, ok := typeCheck[T](m, key, value)
	if !ok {
		return false
	}
	return val != value
}

// unEqualStr !=
func unEqualStr(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return !strings.EqualFold(val, value.(string))
}

// gte >=
func gte[T int64 | float64 | string](m map[string]any, key string, value any) bool {
	val, ok := typeCheck[T](m, key, value)
	if !ok {
		return false
	}
	return val >= value.(T)
}

// lte <=
func lte[T int64 | float64 | string](m map[string]any, key string, value any) bool {
	val, ok := typeCheck[T](m, key, value)
	if !ok {
		return false
	}
	return val <= value.(T)
}

// gt >
func gt[T int64 | float64 | string](m map[string]any, key string, value any) bool {
	val, ok := typeCheck[T](m, key, value)
	if !ok {
		return false
	}
	return val > value.(T)
}

// lt <
func lt[T int64 | float64 | string](m map[string]any, key string, value any) bool {
	val, ok := typeCheck[T](m, key, value)
	if !ok {
		return false
	}
	return val < value.(T)
}

// contains 包含
func contains(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return strings.Contains(strings.ToLower(val), strings.ToLower(value.(string)))
}

// unContains 不包含
func unContains(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return !strings.Contains(strings.ToLower(val), strings.ToLower(value.(string)))
}

// startsWith 前缀匹配
func startsWith(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return strings.HasPrefix(strings.ToLower(val), strings.ToLower(value.(string)))
}

// unStartsWith 前缀不匹配
func unStartsWith(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return !strings.HasPrefix(strings.ToLower(val), strings.ToLower(value.(string)))
}

// endsWith 后缀匹配
func endsWith(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return strings.HasSuffix(strings.ToLower(val), strings.ToLower(value.(string)))
}

// unEndsWith 后缀不匹配
func unEndsWith(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return !strings.HasSuffix(strings.ToLower(val), strings.ToLower(value.(string)))
}

// reg 正则
func reg(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	result, err := regexp.MatchString(value.(string), val)
	if err != nil {
		return false
	}
	return result
}

func in(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return utils.FindIndex(strings.Split(value.(string), ","), val) != -1
}

func notIn(m map[string]any, key string, value any) bool {
	val, ok := typeCheck[string](m, key, value)
	if !ok {
		return false
	}
	return utils.FindIndex(strings.Split(value.(string), ","), val) == -1
}

// containsBit 位运算不进行类型判断，直接转成 int64
func containsBit(m map[string]any, key string, value any) bool {
	mVal, err := mapstr.M(m).GetValue(key)
	if err != nil {
		return false
	}
	mIntVal := number.ParseInt[int64](mVal)
	intVal := number.ParseInt[int64](value)
	return mIntVal&intVal == intVal
}

// containsBit 位运算不进行类型判断，直接转成 int64
func unContainsBit(m map[string]any, key string, value any) bool {
	mVal, err := mapstr.M(m).GetValue(key)
	if err != nil {
		return false
	}
	mIntVal := number.ParseInt[int64](mVal)
	intVal := number.ParseInt[int64](value)
	return mIntVal&intVal != intVal
}
