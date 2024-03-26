package es_expr

import (
	"fmt"

	"github.com/jummyliu/pkg/expression"
	"github.com/jummyliu/pkg/expression/token"
)

type Executor struct {
	FnMap map[string]map[token.Token]ConditionFn
}

var StdExecutor = New(nil)

func New(fnMap map[string]map[token.Token]ConditionFn) *Executor {
	if fnMap == nil {
		fnMap = DefaultFnMap
	}
	return &Executor{
		FnMap: fnMap,
	}
}

// DoExpr 执行表达式
func (e *Executor) DoExpr(expr string, prefix, suffix string) (query map[string]any, keys []string, err error) {
	lexTokens, err := expression.LexParse(expression.TokensRead(expr))
	if err != nil {
		return nil, nil, err
	}
	ast := expression.AstParse(lexTokens)
	query, keys = e.DoAst(ast, prefix, suffix)
	return query, keys, nil
}

// DoAst 执行 ast
func (e *Executor) DoAst(ast *expression.AstNode, prefix, suffix string) (query map[string]any, keys []string) {
	if ast == nil {
		return nil, nil
	}
	switch ast.Type {
	case token.OPERATOR:
		// 复合转换
		leftResult, leftKeys := e.DoAst(ast.Left, prefix, suffix)
		rightResult, rightKeys := e.DoAst(ast.Right, prefix, suffix)
		if leftResult == nil || rightResult == nil {
			return nil, nil
		}
		keys = append(keys, leftKeys...)
		keys = append(keys, rightKeys...)
		switch ast.Value.(string) {
		case "&&":
			return map[string]any{
				"bool": map[string]any{
					"must": []map[string]any{
						leftResult,
						rightResult,
					},
				},
			}, keys
		case "||":
			return map[string]any{
				"bool": map[string]any{
					"should": []map[string]any{
						leftResult,
						rightResult,
					},
				},
			}, keys
		}
	case token.CONDITION:
		// 单个表达式转换
		return e.DoTerm(ast, prefix, suffix)
	}
	return nil, nil
}

// DoTerm 执行 term
func (e *Executor) DoTerm(term *expression.AstNode, prefix, suffix string) (query map[string]any, keys []string) {
	// 判断是否有 left 和 right
	if term.Left == nil || term.Right == nil {
		return nil, nil
	}
	// 判断 left，即 key 的类型
	if _, ok := term.Left.Value.(string); !ok {
		return nil, nil
	}
	// 判断 condition 的类型
	if _, ok := term.Value.(string); !ok {
		return nil, nil
	}
	fns, ok := e.FnMap[term.Value.(string)]
	if !ok {
		return nil, nil
	}
	fn, ok := fns[term.Right.Type]
	if !ok {
		return nil, nil
	}
	key := term.Left.Value.(string)
	if len(prefix) > 0 {
		key = fmt.Sprintf("%s.%s", prefix, key)
	}
	if len(suffix) > 0 {
		key = fmt.Sprintf("%s.%s", key, suffix)
	}
	return fn(key, term.Right.Value), []string{key}
}
