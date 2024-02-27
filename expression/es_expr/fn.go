package es_expr

import (
	"fmt"

	"github.com/jummyliu/pkg/expression/token"
)

type conditionFn func(key string, value any) map[string]any

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
	"&": {},
	"|": {},
}

func equal[T comparable](key string, value any) map[string]any {
	val, ok := value.(T)
	if !ok {
		return nil
	}
	return map[string]any{
		"term": map[string]any{
			key: map[string]any{
				// "case_insensitive": true,
				"value": val,
			},
		},
	}
}

func equalStr(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"term": map[string]any{
			key: map[string]any{
				"case_insensitive": true,
				"value":            val,
			},
		},
	}
}

func unEqual[T comparable](key string, value any) map[string]any {
	val, ok := value.(T)
	if !ok {
		return nil
	}
	return map[string]any{
		"bool": map[string]any{
			"must_not": []map[string]any{
				{
					"term": map[string]any{
						key: map[string]any{
							// "case_insensitive": true,
							"value": val,
						},
					},
				},
			},
		},
	}
}

func unEqualStr(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"bool": map[string]any{
			"must_not": []map[string]any{
				{
					"term": map[string]any{
						key: map[string]any{
							"case_insensitive": true,
							"value":            val,
						},
					},
				},
			},
		},
	}
}

func gte[T int64 | float64 | string](key string, value any) map[string]any {
	val, ok := value.(T)
	if !ok {
		return nil
	}
	return map[string]any{
		"range": map[string]any{
			key: map[string]any{
				"gte": val,
			},
		},
	}
}

func lte[T int64 | float64 | string](key string, value any) map[string]any {
	val, ok := value.(T)
	if !ok {
		return nil
	}
	return map[string]any{
		"range": map[string]any{
			key: map[string]any{
				"lte": val,
			},
		},
	}
}

func gt[T int64 | float64 | string](key string, value any) map[string]any {
	val, ok := value.(T)
	if !ok {
		return nil
	}
	return map[string]any{
		"range": map[string]any{
			key: map[string]any{
				"gt": val,
			},
		},
	}
}

func lt[T int64 | float64 | string](key string, value any) map[string]any {
	val, ok := value.(T)
	if !ok {
		return nil
	}
	return map[string]any{
		"range": map[string]any{
			key: map[string]any{
				"lt": val,
			},
		},
	}
}

func contains(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"wildcard": map[string]any{
			key: map[string]any{
				"case_insensitive": true,
				"value":            fmt.Sprintf("*%s*", val),
			},
		},
	}
}

func unContains(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"bool": map[string]any{
			"must_not": []map[string]any{
				{
					"wildcard": map[string]any{
						key: map[string]any{
							"case_insensitive": true,
							"value":            fmt.Sprintf("*%s*", val),
						},
					},
				},
			},
		},
	}
}

func startsWith(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"wildcard": map[string]any{
			key: map[string]any{
				"case_insensitive": true,
				"value":            fmt.Sprintf("%s*", val),
			},
		},
	}
}

func unStartsWith(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"bool": map[string]any{
			"must_not": []map[string]any{
				{
					"wildcard": map[string]any{
						key: map[string]any{
							"case_insensitive": true,
							"value":            fmt.Sprintf("%s*", val),
						},
					},
				},
			},
		},
	}
}

func endsWith(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"wildcard": map[string]any{
			key: map[string]any{
				"case_insensitive": true,
				"value":            fmt.Sprintf("*%s", val),
			},
		},
	}
}

func unEndsWith(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"bool": map[string]any{
			"must_not": []map[string]any{
				{
					"wildcard": map[string]any{
						key: map[string]any{
							"case_insensitive": true,
							"value":            fmt.Sprintf("*%s", val),
						},
					},
				},
			},
		},
	}
}

func reg(key string, value any) map[string]any {
	val, ok := value.(string)
	if !ok {
		return nil
	}
	key = fmt.Sprintf("%s.keyword", key)
	return map[string]any{
		"regexp": map[string]any{
			key: val,
		},
	}
}
