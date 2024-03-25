package mysql_expr

import (
	"fmt"

	"github.com/jummyliu/pkg/expression/token"
	"github.com/jummyliu/pkg/number"
)

type conditionFn func(key string, value any) (sqls string, params []any)

var DefaultFnMap = map[string]map[token.Token]conditionFn{
	"==": {
		token.NUM:    equal[float64],
		token.BOOL:   equal[bool],
		token.STRING: equal[string],
	},
	"!=": {
		token.NUM:    unEqual[float64],
		token.BOOL:   unEqual[bool],
		token.STRING: unEqual[string],
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

func equal[T comparable](key string, value any) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s = ?", key), []any{val}
}

func unEqual[T comparable](key string, value any) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s != ?", key), []any{val}
}

func gte[T int64 | float64 | string](key string, value any) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s >= ?", key), []any{val}
}

func lte[T int64 | float64 | string](key string, value any) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s <= ?", key), []any{val}
}

func gt[T int64 | float64 | string](key string, value any) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s > ?", key), []any{val}
}

func lt[T int64 | float64 | string](key string, value any) (sql string, params []any) {
	val, ok := value.(T)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s < ?", key), []any{val}
}

func contains(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s LIKE CONCAT('%%', ?, '%%')", key), []any{val}
}

func unContains(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s NOT LIKE CONCAT('%%', ?, '%%')", key), []any{val}
}

func startsWith(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s LIKE CONCAT(?, '%%')", key), []any{val}
}

func unStartsWith(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s NOT LIKE CONCAT(?, '%%')", key), []any{val}
}

func endsWith(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s LIKE CONCAT('%%', ?)", key), []any{val}
}

func unEndsWith(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s NOT LIKE CONCAT('%%', ?)", key), []any{val}
}

func reg(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("%s REGEXP ?", key), []any{val}
}

func in(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("FIND_IN_SET(`%s`, ?)", key), []any{val}
}

func notIn(key string, value any) (sql string, params []any) {
	val, ok := value.(string)
	if !ok {
		return "", nil
	}
	return fmt.Sprintf("NOT FIND_IN_SET(`%s`, ?)", key), []any{val}
}

// containsBit 位运算不进行类型判断，直接转成 int64
func containsBit(key string, value any) (sql string, params []any) {
	intVal := number.ParseInt[int64](value)
	return fmt.Sprintf("%s & ? = ?", key), []any{intVal, intVal}
}

// containsBit 位运算不进行类型判断，直接转成 int64
func unContainsBit(key string, value any) (sql string, params []any) {
	intVal := number.ParseInt[int64](value)
	return fmt.Sprintf("%s & ? != ?", key), []any{intVal, intVal}
}
