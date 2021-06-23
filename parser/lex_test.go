package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	testcases := []struct {
		content string
		tokens  []Token
	}{
		{
			content: "{()}",
			tokens: []Token{
				LEFT_BRACE.toToken(),
				LEFT_PAREN.toToken(),
				RIGHT_PAREN.toToken(),
				RIGHT_BRACE.toToken(),
			},
		},
		{
			content: `
			// Hello world
			123.12 123.
			`,
			tokens: []Token{
				{TokenType: NEWLINE, lexeme: "", line: 0},
				{TokenType: WHITESPACE, lexeme: "\t\t", line: 0},
				{TokenType: LINE_COMMENT, lexeme: " Hello world", line: 0},
				{TokenType: NEWLINE, lexeme: "", line: 0},
				{TokenType: WHITESPACE, lexeme: "\t\t", line: 0},
				{TokenType: NUMBER, lexeme: "123.12", line: 0},
				{TokenType: WHITESPACE, lexeme: "", line: 0},
				{TokenType: NUMBER, lexeme: "123", line: 0},
				{TokenType: DOT, lexeme: "", line: 0},
				{TokenType: NEWLINE, lexeme: "", line: 0},
				{TokenType: WHITESPACE, lexeme: "\t\t", line: 0}},
		},
	}

	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			tokens, err := Lex(tc.content)
			require.NoError(t, err)
			assert.Equal(t, tc.tokens, tokens)
		})
	}
}
