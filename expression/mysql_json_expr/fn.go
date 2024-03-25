package mysql_json_expr

import (
	"fmt"
	"strings"

	"github.com/jummyliu/pkg/expression/token"
	"github.com/jummyliu/pkg/number"
)

func buildKey(key string) (sql string, params []any) {
	if len(key) == 0 {
		return "", nil
	}
	keys := strings.Split(key, ".")
	arr := make([]string, 0, len(keys))
	params = make([]any, 0, len(keys))
	for _, k := range keys {
		arr = append(arr, "'\"', ?, '\"'")
		params = append(params, k)
	}
	return fmt.Sprintf("CONCAT('$.', %s, '')", strings.Join(arr, ", '.', ")), params
}

type conditionFn func(key string, value any, jsonAttr string) (sqls string, params []any)

var DefaultFnMap = map[string]map[token.Token]conditionFn{
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

func equal[T comparable](key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_CONTAINS(%s, CONCAT('', ?, ''), %s)",
		jsonAttr,
		keySql,
	), params
}

func equalStr(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_CONTAINS(%s, CONCAT('\"', ?, '\"'), %s)",
		jsonAttr,
		keySql,
	), params
}

func unEqual[T comparable](key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_CONTAINS(%s, CONCAT('', ?, ''), %s) = 0",
		jsonAttr,
		keySql,
	), params
}

func unEqualStr(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_CONTAINS(%s, CONCAT('\"', ?, '\"'), %s) = 0",
		jsonAttr,
		keySql,
	), params
}

func gte[T int64 | float64 | string](key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, val)
	return fmt.Sprintf(
		"JSON_EXTRACT(%s, %s) COLLATE utf8mb4_0900_ai_ci >= ?",
		jsonAttr,
		keySql,
	), params
}

func lte[T int64 | float64 | string](key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, val)
	return fmt.Sprintf(
		"JSON_EXTRACT(%s, %s) COLLATE utf8mb4_0900_ai_ci <= ?",
		jsonAttr,
		keySql,
	), params
}

func gt[T int64 | float64 | string](key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, val)
	return fmt.Sprintf(
		"JSON_EXTRACT(%s, %s) COLLATE utf8mb4_0900_ai_ci > ?",
		jsonAttr,
		keySql,
	), params
}

func lt[T int64 | float64 | string](key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, val)
	return fmt.Sprintf(
		"JSON_EXTRACT(%s, %s) COLLATE utf8mb4_0900_ai_ci < ?",
		jsonAttr,
		keySql,
	), params
}

func contains(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_SEARCH(%s, 'one', CONCAT('%%', ?, '%%'), null, %s) IS NOT NULL",
		jsonAttr,
		keySql,
	), params
}

func unContains(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_SEARCH(%s, 'one', CONCAT('%%', ?, '%%'), null, %s) IS NULL",
		jsonAttr,
		keySql,
	), params
}

func startsWith(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_SEARCH(%s, 'one', CONCAT(?, '%%'), null, %s) IS NOT NULL",
		jsonAttr,
		keySql,
	), params
}

func unStartsWith(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_SEARCH(%s, 'one', CONCAT(?, '%%'), null, %s) IS NULL",
		jsonAttr,
		keySql,
	), params
}

func endsWith(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_SEARCH(%s, 'one', CONCAT('%%', ?), null, %s) IS NOT NULL",
		jsonAttr,
		keySql,
	), params
}

func unEndsWith(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, val)
	params = append(params, p...)
	return fmt.Sprintf(
		"JSON_SEARCH(%s, 'one', CONCAT('%%', ?), null, %s) IS NULL",
		jsonAttr,
		keySql,
	), params
}

func reg(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, val)
	return fmt.Sprintf(
		"JSON_EXTRACT(%s, %s) COLLATE utf8mb4_0900_ai_ci REGEXP ?",
		jsonAttr,
		keySql,
	), params
}

func in(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, val)
	return fmt.Sprintf(
		"FIND_IN_SET(JSON_UNQUOTE(JSON_EXTRACT(%s , %s)), ?)",
		jsonAttr,
		keySql,
	), params
}

func notIn(key string, value any, jsonAttr string) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, val)
	return fmt.Sprintf(
		"!FIND_IN_SET(JSON_UNQUOTE(JSON_EXTRACT(%s , %s)), ?)",
		jsonAttr,
		keySql,
	), params
}

// containsBit 位运算不进行类型判断，直接转成 int64
func containsBit(key string, value any, jsonAttr string) (sql string, params []any) {
	intVal := number.ParseInt[int64](value)
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, intVal, intVal)
	return fmt.Sprintf(
		"JSON_EXTRACT(%s, %s) & ? = ?",
		jsonAttr,
		keySql,
	), params
}

// containsBit 位运算不进行类型判断，直接转成 int64
func unContainsBit(key string, value any, jsonAttr string) (sql string, params []any) {
	intVal := number.ParseInt[int64](value)
	keySql, p := buildKey(key)
	params = append(params, p...)
	params = append(params, intVal, intVal)
	return fmt.Sprintf(
		"JSON_EXTRACT(%s, %s) & ? != ?",
		jsonAttr,
		keySql,
	), params
}
