package token

import (
	"regexp"
)

// token 的正则表达式
var (
	REG_LBT       = regexp.MustCompile(`^([(\[{])`)
	REG_RBT       = regexp.MustCompile(`^([)\]}])`)
	REG_OPERATOR  = regexp.MustCompile(`^(&&|\|\|)`)
	REG_CONDITION = regexp.MustCompile(`^(?:(==|!=|>=|<=|>|<|containsBit|unContainsBit|contains|unContains|startsWith|unStartsWith|endsWith|unEndsWith|reg)\s*)`)
	REG_IDENT     = regexp.MustCompile(`^([a-zA-Z_][\w\\.\-]*)`)
	REG_NUM       = regexp.MustCompile(`^((?:\+|-)?\d+(?:\.\d+)?)`)
	REG_BOOL      = regexp.MustCompile(`^(true|false)`)
	REG_STRING    = regexp.MustCompile(`^(?:'(.*?)'|"(.*?)")`)
	REG_DELIM     = regexp.MustCompile(`^(\s*)`)
	REG_ILLEGAL   = regexp.MustCompile(`^(.+)`)
)

type Token string

// token
const (
	LITERAL_BEGIN Token = "literal_begin"
	LITERAL_END   Token = "literal_end"

	LBT       Token = "lbt"       // 左括号
	RBT       Token = "rbt"       // 有括号
	OPERATOR  Token = "operator"  // 逻辑关系
	CONDITION Token = "condition" // 条件
	IDENT     Token = "ident"     // 变量
	NUM       Token = "num"       // 数字
	BOOL      Token = "bool"      // bool
	STRING    Token = "string"    // 字符串
	DELIM     Token = "delim"     // 空字符
	ILLEGAL   Token = "illegal"   // 无法解析
)

// Parser token 解析器
//
//	包含词法分析正则
//	包含 Token 标识
//	包含编码、解码器
type Parser struct {
	Reg   *regexp.Regexp
	Token Token
	Encodable
}

// token 解析器
var (
	lbtParser       = Parser{Reg: REG_LBT, Token: LBT, Encodable: EmptyEncoder}
	rbtParser       = Parser{Reg: REG_RBT, Token: RBT, Encodable: EmptyEncoder}
	operatorParser  = Parser{Reg: REG_OPERATOR, Token: OPERATOR, Encodable: EmptyEncoder}
	conditionParser = Parser{Reg: REG_CONDITION, Token: CONDITION, Encodable: EmptyEncoder}
	identParser     = Parser{Reg: REG_IDENT, Token: IDENT, Encodable: EmptyEncoder}
	numParser       = Parser{Reg: REG_NUM, Token: NUM, Encodable: NumberEncoder}
	boolParser      = Parser{Reg: REG_BOOL, Token: BOOL, Encodable: BoolEncoder}
	stringParser    = Parser{Reg: REG_STRING, Token: STRING, Encodable: StringEncoder}
	DelimParser     = Parser{Reg: REG_DELIM, Token: DELIM, Encodable: EmptyEncoder}
	illegalParser   = Parser{Reg: REG_ILLEGAL, Token: ILLEGAL, Encodable: EmptyEncoder}
)

// ParserPriority token 解析器优先级
var ParserPriority = []Parser{
	conditionParser,
	operatorParser,
	numParser,
	boolParser,
	stringParser,
	identParser,
	lbtParser,
	rbtParser,
	illegalParser,
}

// ParserMap token 解析器映射
var ParserMap = map[Token]Parser{
	LBT:       lbtParser,
	RBT:       rbtParser,
	OPERATOR:  operatorParser,
	CONDITION: conditionParser,
	IDENT:     identParser,
	NUM:       numParser,
	BOOL:      boolParser,
	STRING:    stringParser,
	DELIM:     DelimParser,
	ILLEGAL:   illegalParser,
}

// StateMatrix 状态转移矩阵
var StateMatrix = map[Token][]Token{
	LITERAL_BEGIN: {LBT, IDENT},
	LBT:           {LBT, IDENT},
	RBT:           {OPERATOR, RBT, LITERAL_END},
	OPERATOR:      {LBT, IDENT},
	CONDITION:     {NUM, BOOL, STRING},
	IDENT:         {CONDITION},
	NUM:           {OPERATOR, RBT, LITERAL_END},
	BOOL:          {OPERATOR, RBT, LITERAL_END},
	STRING:        {OPERATOR, RBT, LITERAL_END},
	ILLEGAL:       {LITERAL_END},
}
