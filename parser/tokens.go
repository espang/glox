package parser

import "fmt"

type TokenType int

func (tt TokenType) toToken() Token {
	return Token{
		TokenType: tt,
	}
}

func (tt TokenType) String() string {
	return tokenTypeStrings(tt)
}

const (
	ILLEGAL TokenType = iota
	EOF
	LINE_COMMENT
	WHITESPACE
	NEWLINE

	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	NOT
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	IDENTIFIER
	STRING
	NUMBER

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
)

type Token struct {
	TokenType

	lexeme string
	line   int
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, content: %s, line: %d)", t.TokenType, t.lexeme, t.line)
}

func keywords() map[string]TokenType {
	return map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
}

// if actually just for more than error logging/ repl this should be made different
func tokenTypeStrings(tt TokenType) string {
	m := map[TokenType]string{
		ILLEGAL:       "illegal",
		EOF:           "EOF",
		LINE_COMMENT:  "line comment",
		WHITESPACE:    "whitespace",
		NEWLINE:       "new line",
		LEFT_PAREN:    "left parentheses",
		RIGHT_PAREN:   "right parentheses",
		LEFT_BRACE:    "left brackets",
		RIGHT_BRACE:   "right brackets",
		COMMA:         "comma",
		DOT:           "dot",
		MINUS:         "minus",
		PLUS:          "plus",
		SEMICOLON:     "semicolon",
		SLASH:         "slash",
		STAR:          "star",
		NOT:           "not",
		BANG_EQUAL:    "not equal",
		EQUAL:         "equal",
		EQUAL_EQUAL:   "equal equal",
		GREATER:       "greater",
		GREATER_EQUAL: "greater equal",
		LESS:          "less",
		LESS_EQUAL:    "less equal",
		IDENTIFIER:    "identifier",
		STRING:        "string",
		NUMBER:        "number",
		AND:           "and",
		CLASS:         "class",
		ELSE:          "else",
		FALSE:         "false",
		FUN:           "function",
		FOR:           "for",
		IF:            "if",
		NIL:           "nil",
		OR:            "or",
		PRINT:         "print",
		RETURN:        "return",
		SUPER:         "super",
		THIS:          "this",
		TRUE:          "true",
		VAR:           "var",
		WHILE:         "while",
	}
	return m[tt]
}
