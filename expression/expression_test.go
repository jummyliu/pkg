package expression

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTokenRead(t *testing.T) {
	testCases := []struct {
		Expr   string
		Result []*LexNode
	}{
		{
			Expr: "test == 123 && ( keyword contains 'hello'))",
			Result: []*LexNode{
				{Type: "ident", Value: "test", Len: 4, From: 0},
				{Type: "condition", Value: "==", Len: 3, From: 5},
				{Type: "num", Value: float64(123), Len: 3, From: 8},
				{Type: "operator", Value: "&&", Len: 2, From: 12},
				{Type: "lbt", Value: "(", Len: 1, From: 15},
				{Type: "ident", Value: "keyword", Len: 7, From: 17},
				{Type: "condition", Value: "contains", Len: 9, From: 25},
				{Type: "string", Value: "hello", Len: 7, From: 34},
				{Type: "rbt", Value: ")", Len: 1, From: 41},
				{Type: "rbt", Value: ")", Len: 1, From: 42},
			},
		},
	}
	for _, testCase := range testCases {
		results := TokensRead(testCase.Expr)
		if !compareTokenObj(testCase.Result, results) {
			t.Logf("testCase %s need:\n", testCase.Expr)
			PrintTokenObj(t, testCase.Result)
			t.Logf("but got:\n")
			PrintTokenObj(t, results)
			t.Fail()
		}
	}
}

func TestLexParse(t *testing.T) {
	testCases := []struct {
		Expr   string
		Result []*LexNode
	}{
		{
			Expr: "test == 123 && ( keyword contains 'hello' && a unContainsBit 10)",
		},
	}
	for _, testCase := range testCases {
		results, err := LexParse(TokensRead(testCase.Expr))
		PrintTokenObj(t, results)
		fmt.Println(err)
		if err != nil {
			t.Fail()
		}
		expr := LexToExpr(results)
		fmt.Println(expr)
		// t.Fail()
	}
}

func TestAstParse(t *testing.T) {
	testCases := []struct {
		Expr   string
		Result []*LexNode
	}{
		{
			Expr: "aaa == 10 && hello != true || (_term >= '2012-12-22' && _term <= '2012-01-01') && asdf != true || abe == 10",
		},
	}
	for _, testCase := range testCases {
		results, err := LexParse(TokensRead(testCase.Expr))
		if err != nil {
			PrintTokenObj(t, results)
			fmt.Println(err)
			t.FailNow()
		}
		ast := AstParse(results)
		data, err := json.Marshal(ast)
		fmt.Println(string(data), err)
		// t.Fail()
	}
}

func compareTokenObj(from, to []*LexNode) bool {
	if from == nil && to == nil {
		return true
	}
	if len(from) != len(to) {
		return false
	}
	for i := range from {
		if !(from[i].From == to[i].From &&
			from[i].Len == to[i].Len &&
			from[i].Type == to[i].Type &&
			from[i].Value == to[i].Value) {
			return false
		}
	}
	return true
}

func PrintTokenObj(t *testing.T, obj []*LexNode) {
	for _, item := range obj {
		if len(item.SubCond) != 0 {
			t.Log(">>>")
			PrintTokenObj(t, item.SubCond)
			t.Log("<<<")
		} else {
			t.Logf("%#v\n", item)
		}
	}
}
