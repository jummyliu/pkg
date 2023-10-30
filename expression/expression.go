// 表达式文法分析
// expr       => expr op term | term
//      => 消除左递归
//      expr  => term expr1
//      expr1 => op term expr1 | null
// term       => key cond val | '(' expr ')'
// op         => '&&' | '||'
//            => /(&&|\|\|)/
// cond       => '==' | '!=' | '>=' | '<=' | '>' | '<'
//            => /(==|!=|>=|<=|>|<)/     // >= <= 要在 > < 前面；不然会匹配不到
// key        => /[\w_][\w\d_]*/
// val        => number | bool | string
//            // 正负数字
//            => /(?:\+|-)?\d+(?:\.\d+)?/
//            // true | false
//            => /(true|false)/
//            // 首尾匹配 ' | "
//            // (?<!\\)' 负后顾，'前面不是单个斜杠
//            // (?<=(?<!\\)(?:\\\\)+)' 后顾，'前面必须是偶数个斜杠
//            => /^('|")(.*?)(?:(?<!\\)|(?<=(?<!\\)(?:\\\\)+))(\1)/
/* ---------------------------------------------------------------- */

/* ---------------------------------------------------------------- */
// 0. 词法分析，把所有词转换成 token(type, val)
//
// 括号
//    LBT   	/[(\[{]/
//    RBT   	/[)\]}]/
//
// 关键字
//    // KEYWORD  //
//
// 逻辑
//    OP    	/(&&|\|\|)/
//
// 条件
//    COND  	/(==|!=|>=|<=|>|<)/
//
// 变量
//    IDENT 	/[a-zA-Z_][\w]*/
//
// 值
//    VAL   NUM | BOOL | STRING
//
//    NUM   	/(?:\+|-)?\d+(?:\.\d+)?/
//    BOOL  	/(?:true|false)/
//    STRING	/('|")(.*?)(?:(?<!\\)|(?<=(?<!\\)(?:\\\\)+))(\1)/
//
// 空字符
//    DELIM		/\s*/
//
// 无法解析
//    ILLEGAL	/.+/
//
// 优先级 COND = OP > VAL > IDENT > LBT = RBT > (DELIM: 每次匹配前过滤掉空字符) > ILLEGAL
/* ---------------------------------------------------------------- */
/* ---------------------------------------------------------------- */
// 1. 把分词进行解析，生成最终的词法数组 >> 可以转回表达式
/* ---------------------------------------------------------------- */
/* ---------------------------------------------------------------- */
// 2. 把词法数组解析，转成语法树 >> 不可以转回表达式
package expression

import (
	"errors"
	"strings"

	"github.com/jummyliu/pkg/expression/token"
	"github.com/jummyliu/pkg/utils"
)

// LexNode 分词节点
type LexNode struct {
	Type    token.Token
	Value   any
	SubCond []*LexNode // 不为空，则为子表达式
	Len     int
	From    int
}

// TokensRead 读取表达式所有分词
//
//	正常来说，返回值是个中间状态，一般不直接使用，需要调用 TokensParse 进行解析
func TokensRead(expr string) (tokens []*LexNode) {
	var index = 0
	for len(expr) != 0 {
		tok := next(expr, token.DelimParser)
		if tok != nil {
			index += tok.Len
			expr = expr[tok.Len:]
		}
		for _, meta := range token.ParserPriority {
			tok := next(expr, meta)
			if tok != nil {
				tok.From = index
				tokens = append(tokens, tok)
				index += tok.Len
				expr = expr[tok.Len:]
				break
			}
		}
	}
	return tokens
}

// 获取下一个 token
func next(expr string, tokenParser token.Parser) *LexNode {
	result := tokenParser.Reg.FindString(expr)
	if len(result) == 0 {
		return nil
	}
	return &LexNode{
		Type:  tokenParser.Token,
		Value: tokenParser.Decode(strings.TrimRight(result, " ")),
		Len:   len(result),
	}
}

// LexParse 解析所有的分词
func LexParse(tokens []*LexNode) (lexTokens []*LexNode, err error) {
	// 为空
	if len(tokens) == 0 {
		return nil, nil
	}
	newTokens := []*LexNode{}
	stack := []*[]*LexNode{}
	curState := token.LITERAL_BEGIN
	var active *[]*LexNode = &newTokens
	last := 0
	for i, item := range tokens {
		if utils.FindIndex[token.Token](token.StateMatrix[curState], item.Type) == -1 {
			// 异常，下一状态不匹配
			newTokens = append(newTokens, &LexNode{
				Type:  token.ILLEGAL,
				Value: errors.New("illegal NEXT STATE"),
			})
			break
		}
		if item.Type == token.LBT {
			// 左括号，把数据压入栈
			tmp := &LexNode{
				SubCond: []*LexNode{},
			}
			stack = append(stack, active)
			*active = append(*active, tmp)
			active = &tmp.SubCond
			*active = append(*active, item)
		} else if item.Type == token.RBT {
			// 右括号，从栈里取最后一个数据
			if len(stack) == 0 {
				// 异常，右括号不匹配
				newTokens = append(newTokens, &LexNode{
					Type:  token.ILLEGAL,
					Value: errors.New("illegal MISS LBT"),
				})
				break
			}
			tmp := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			*active = append(*active, item)
			active = tmp
		} else {
			*active = append(*active, item)
		}
		curState = item.Type
		last = i + 1
	}
	if newTokens[len(newTokens)-1].Type != token.ILLEGAL {
		// 判断是否解析完成
		// 最后状态没有 end
		if last != len(tokens) {
			// 没有解析完 token，异常
			newTokens = append(newTokens, &LexNode{
				Type:  token.ILLEGAL,
				Value: errors.New("illegal EOF"),
			})
		} else if utils.FindIndex[token.Token](token.StateMatrix[curState], token.LITERAL_END) == -1 {
			newTokens = append(newTokens, &LexNode{
				Type:  token.ILLEGAL,
				Value: errors.New("illegal MISS END"),
			})
		} else if len(stack) != 0 {
			//  括号没闭合
			newTokens = append(newTokens, &LexNode{
				Type:  token.ILLEGAL,
				Value: errors.New("illegal MISS RBT"),
			})
		}
	}
	// 解析出现了错误，返回异常
	if newTokens[len(newTokens)-1].Type == token.ILLEGAL {
		return newTokens, newTokens[len(newTokens)-1].Value.(error)
	}
	return newTokens, nil
}

// LexToExpr 分词转表达式
func LexToExpr(tokens []*LexNode) (expr string) {
	if len(tokens) == 0 {
		return
	}
	var b strings.Builder
	b.Grow(len(tokens))
	for _, item := range tokens {
		if len(item.SubCond) != 0 {
			b.WriteByte(' ')
			b.WriteString(LexToExpr(item.SubCond))
			continue
		}
		parser, ok := token.ParserMap[item.Type]
		if !ok {
			return ""
		}
		b.WriteByte(' ')
		b.WriteString(parser.Encode(item.Value))
	}
	return b.String()[1:]
}

// AstNode 语法树节点
type AstNode struct {
	Type  token.Token `json:"type"`
	Value any         `json:"value"`
	Left  *AstNode    `json:"left"`  // sub left tree, 只有非 term 节点才有, term 节点为 nil
	Right *AstNode    `json:"right"` // sub right tree, 只有非 term 节点有, term 节点为 nil
}

// AstParse 把分词解析成语法树
func AstParse(lexTokens []*LexNode) (tree *AstNode) {
	var root *AstNode
	var node *AstNode
	for i := 0; i < len(lexTokens); i++ {
		item := lexTokens[i]
		if len(item.SubCond) != 0 {
			node = AstParse(item.SubCond)
			if root == nil {
				root = node
			} else {
				root.Right = node
			}
			continue
		}
		switch item.Type {
		case token.OPERATOR:
			if root == nil {
				break
			}
			node = &AstNode{
				Type:  item.Type,
				Value: item.Value,
				Left:  root,
			}
			root = node
		case token.IDENT:
			condition := lexTokens[i+1]
			value := lexTokens[i+2]
			node = &AstNode{
				Type:  condition.Type,
				Value: condition.Value,
				Left: &AstNode{
					Type:  item.Type,
					Value: item.Value,
				},
				Right: &AstNode{
					Type:  value.Type,
					Value: value.Value,
				},
			}
			if root == nil {
				root = node
			} else {
				root.Right = node
			}
			i += 2
		default:
		}

	}
	return root
}
