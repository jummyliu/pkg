package clickhouse_expr

import (
	"fmt"
	"strings"

	"github.com/jummyliu/pkg/expression"
	"github.com/jummyliu/pkg/expression/token"
)

const (
	lbt         = "( "
	rbt         = " )"
	operatorAnd = " AND "
	operatorOr  = " OR "

	lenAnd = len(lbt) + len(rbt) + len(operatorAnd)
	lenOr  = len(lbt) + len(rbt) + len(operatorOr)
)

type Executor struct {
	FnMap map[string]map[token.Token]ConditionFn

	// KeyMap 字段映射
	// 	存在映射 => key 转换为映射值
	KeyMap map[string]string
}

var StdExecutor = New(nil, nil)

func New(fnMap map[string]map[token.Token]ConditionFn, keyMap map[string]string) *Executor {
	if fnMap == nil {
		fnMap = DefaultFnMap
	}
	return &Executor{
		FnMap: fnMap,

		KeyMap: keyMap,
	}
}

// DoExpr 执行表达式
func (e *Executor) DoExpr(expr string, prefix, suffix string) (sqls string, params []any, keys []string, err error) {
	lexTokens, err := expression.LexParse(expression.TokensRead(expr))
	if err != nil {
		return "", nil, nil, err
	}
	ast := expression.AstParse(lexTokens)
	sqls, params, keys = e.DoAst(ast, prefix, suffix)
	return sqls, params, keys, nil
}

// DoAst 执行 ast
func (e *Executor) DoAst(ast *expression.AstNode, prefix, suffix string) (sqls string, params []any, keys []string) {
	if ast == nil {
		return "", nil, nil
	}
	switch ast.Type {
	case token.OPERATOR:
		// 复合转换
		leftSQL, leftParams, leftKeys := e.DoAst(ast.Left, prefix, suffix)
		rightSQL, rightParams, rightKeys := e.DoAst(ast.Right, prefix, suffix)
		if len(leftSQL) == 0 || len(rightSQL) == 0 {
			return "", nil, nil
		}
		switch ast.Value.(string) {
		case "&&":
			var b strings.Builder
			b.Grow(len(leftSQL) + len(rightSQL) + lenAnd)
			b.WriteString(lbt)
			b.WriteString(leftSQL)
			b.WriteString(operatorAnd)
			b.WriteString(rightSQL)
			b.WriteString(rbt)
			params = append(params, leftParams...)
			params = append(params, rightParams...)
			keys = append(keys, leftKeys...)
			keys = append(keys, rightKeys...)
			return b.String(), params, keys
		case "||":
			var b strings.Builder
			b.Grow(len(leftSQL) + len(rightSQL) + lenOr)
			b.WriteString(lbt)
			b.WriteString(leftSQL)
			b.WriteString(operatorOr)
			b.WriteString(rightSQL)
			b.WriteString(rbt)
			params = append(params, leftParams...)
			params = append(params, rightParams...)
			keys = append(keys, leftKeys...)
			keys = append(keys, rightKeys...)
			return b.String(), params, keys
		}
	case token.CONDITION:
		// 单个表达式转换
		return e.DoTerm(ast, prefix, suffix)
	}
	return "", nil, nil
}

// DoTerm 执行 term
func (e *Executor) DoTerm(term *expression.AstNode, prefix, suffix string) (sql string, params []any, keys []string) {
	// 判断是否有 left 和 right
	if term.Left == nil || term.Right == nil {
		return "", nil, nil
	}
	// 判断 left，即 key 的类型
	if _, ok := term.Left.Value.(string); !ok {
		return "", nil, nil
	}
	// 判断 condition 的类型
	if _, ok := term.Value.(string); !ok {
		return "", nil, nil
	}
	fns, ok := e.FnMap[term.Value.(string)]
	if !ok {
		return "", nil, nil
	}
	fn, ok := fns[term.Right.Type]
	if !ok {
		return "", nil, nil
	}
	key := term.Left.Value.(string)
	if _, ok := e.KeyMap[key]; ok {
		key = e.KeyMap[key]
	}
	if len(prefix) > 0 {
		key = fmt.Sprintf("%s.%s", prefix, key)
	}
	if len(suffix) > 0 {
		key = fmt.Sprintf("%s.%s", key, suffix)
	}
	sql, params = fn(key, term.Right.Value)
	return sql, params, []string{key}
}
