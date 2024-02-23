package token

import (
	"testing"
)

func TestToken(t *testing.T) {
	testCases := []struct {
		Token       string
		TokenParser Parser
		Result      any
	}{
		{
			Token:       "(",
			TokenParser: lbtParser,
			Result:      "(",
		},
		{
			Token:       ")",
			TokenParser: rbtParser,
			Result:      ")",
		},
		{
			Token:       "&&",
			TokenParser: operatorParser,
			Result:      "&&",
		},
		{
			Token:       "&",
			TokenParser: conditionParser,
			Result:      "&",
		},
		{
			Token:       "||",
			TokenParser: operatorParser,
			Result:      "||",
		},
		{
			Token:       "==",
			TokenParser: conditionParser,
			Result:      "==",
		},
		{
			Token:       "test",
			TokenParser: identParser,
			Result:      "test",
		},
		{
			Token:       "123.45",
			TokenParser: numParser,
			Result:      123.45,
		},
		{
			Token:       "true",
			TokenParser: boolParser,
			Result:      true,
		},
		{
			Token:       "\"hello world\"",
			TokenParser: stringParser,
			Result:      "hello world",
		},
		{
			Token:       "'hello world'",
			TokenParser: stringParser,
			Result:      "hello world",
		},
		{
			Token:       "asdfc",
			TokenParser: stringParser,
			Result:      "",
		},
	}

	for _, testCase := range testCases {
		result := testCase.TokenParser.Reg.FindString(testCase.Token)
		resultAny := testCase.TokenParser.Decode(result)
		if resultAny != testCase.Result {
			t.Fatalf("testCase %s need (%v) but got (%v)", testCase.TokenParser.Token, testCase.Result, resultAny)
		}
	}
}
