package es_expr

import (
	"fmt"

	"github.com/jummyliu/pkg/expression"
	"github.com/jummyliu/pkg/expression/token"
)

type Executor struct {
	FnMap map[string]map[token.Token]conditionFn
}

var StdExecutor = New(nil)

func New(fnMap map[string]map[token.Token]conditionFn) *Executor {
	if fnMap == nil {
		fnMap = DefaultFnMap
	}
	return &Executor{
		FnMap: fnMap,
	}
}

// DoExpr 执行表达式
func (e *Executor) DoExpr(m map[string]any, expr string, prefix, suffix string) (result bool, keys []string, err error) {
	lexTokens, err := expression.LexParse(expression.TokensRead(expr))
	if err != nil {
		return false, nil, err
	}
	ast := expression.AstParse(lexTokens)
	result, keys = e.DoAst(m, ast, prefix, suffix)
	return result, keys, nil
}

// DoAst 执行 ast
func (e *Executor) DoAst(m map[string]any, ast *expression.AstNode, prefix, suffix string) (result bool, keys []string) {
	if ast == nil {
		return true, nil
	}
	switch ast.Type {
	case token.OPERATOR:
		// 复合转换
		leftResult, leftKeys := e.DoAst(m, ast.Left, prefix, suffix)
		keys = append(keys, leftKeys...)
		switch ast.Value.(string) {
		case "&&":
			if !leftResult {
				// 提前退出逻辑判断
				return false, keys
			}
			rightResult, rightKeys := e.DoAst(m, ast.Right, prefix, suffix)
			keys = append(keys, rightKeys...)
			return OperatorAnd(leftResult, rightResult), keys
		case "||":
			if leftResult {
				// 提前退出逻辑判断
				return true, keys
			}
			rightResult, rightKeys := e.DoAst(m, ast.Right, prefix, suffix)
			keys = append(keys, rightKeys...)
			return OperatorOr(leftResult, rightResult), keys
		}
	case token.CONDITION:
		// 单个表达式转换
		return e.DoTerm(m, ast, prefix, suffix)
	}
	return false, nil
}

// DoTerm 执行 term
func (e *Executor) DoTerm(m map[string]any, term *expression.AstNode, prefix, suffix string) (result bool, keys []string) {
	// 判断是否有 left 和 right
	if term.Left == nil || term.Right == nil {
		return false, nil
	}
	// 判断 left，即 key 的类型
	if _, ok := term.Left.Value.(string); !ok {
		return false, nil
	}
	// 判断 condition 的类型
	if _, ok := term.Value.(string); !ok {
		return false, nil
	}
	fns, ok := e.FnMap[term.Value.(string)]
	if !ok {
		return false, nil
	}
	fn, ok := fns[term.Right.Type]
	if !ok {
		return false, nil
	}
	key := term.Left.Value.(string)
	if len(prefix) > 0 {
		key = fmt.Sprintf("%s.%s", prefix, key)
	}
	if len(suffix) > 0 {
		key = fmt.Sprintf("%s.%s", key, suffix)
	}
	return fn(m, key, term.Right.Value), []string{key}
}

func OperatorAnd(left, right bool) bool {
	return left && right
}

func OperatorOr(left, right bool) bool {
	return left || right
}
