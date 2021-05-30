package parser

import (
	"container/list"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type lexer struct {
	content         *strings.Reader
	readAheadBuffer *list.List
	keywords        map[string]TokenType
	size            int
	start           int
	current         int
	line            int
}

func lex(content string) *lexer {
	return &lexer{
		content:         strings.NewReader(content),
		readAheadBuffer: list.New(),
		keywords:        keywords(),
		size:            len(content),
		start:           0,
		current:         0,
		line:            1,
	}
}

func (l *lexer) scanToken() Token {
	r, ok := l.read()
	if !ok {
		return Token{TokenType: EOF}
	}
	switch r {
	case '(':
		return LEFT_PAREN.toToken()
	case ')':
		return RIGHT_PAREN.toToken()
	case '{':
		return LEFT_BRACE.toToken()
	case '}':
		return RIGHT_BRACE.toToken()
	case ',':
		return COMMA.toToken()
	case '.':
		return DOT.toToken()
	case '-':
		return MINUS.toToken()
	case '+':
		return PLUS.toToken()
	case ';':
		return SEMICOLON.toToken()
	case '*':
		return STAR.toToken()
	case '!':
		next, ok := l.peak()
		if ok && next == '=' {
			_, _ = l.read()
			return BANG_EQUAL.toToken()
		}
		return NOT.toToken()
	case '=':
		next, ok := l.peak()
		if ok && next == '=' {
			_, _ = l.read()
			return EQUAL_EQUAL.toToken()
		}
		return EQUAL.toToken()
	case '<':
		next, ok := l.peak()
		if ok && next == '=' {
			_, _ = l.read()
			return LESS_EQUAL.toToken()
		}
		return LESS.toToken()
	case '>':
		next, ok := l.peak()
		if ok && next == '=' {
			_, _ = l.read()
			return GREATER_EQUAL.toToken()
		}
		return GREATER.toToken()
	case '/':
		next, ok := l.peak()
		if ok && next == '/' {
			_, _ = l.read()
			comment := l.consumeLineComment()
			return Token{
				TokenType: LINE_COMMENT,
				lexeme:    comment,
			}
		}
		return SLASH.toToken()
	case ' ', '\t', '\r':
		ws := l.consumeWhiteSpace()
		return Token{
			TokenType: WHITESPACE,
			lexeme:    ws,
		}
	case '\n':
		l.line++
		return NEWLINE.toToken()
	case '"':
		content, ok := l.consumeString()
		if !ok {
			panic("unterminated string")
		}
		return Token{
			TokenType: STRING,
			lexeme:    content,
		}
	default:
		if isDigit(r) {
			content := l.consumeNumber(r)
			return Token{
				TokenType: NUMBER,
				lexeme:    content,
			}
		}
		if isAlpha(r) {
			content := l.consumeIdentifier(r)
			if _, ok := l.keywords[content]; ok {
				return Token{
					TokenType: l.keywords[content],
					lexeme:    content,
				}
			}
			return Token{
				TokenType: IDENTIFIER,
				lexeme:    content,
			}
		}
		fmt.Printf("unexpected rune '%v'\n", r)
		return ILLEGAL.toToken()
	}
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlpha(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func (l *lexer) consumeDigits() string {
	var runes []rune
	r, ok := l.peak()
	for ok && isDigit(r) {
		_, _ = l.read()
		runes = append(runes, r)
		r, ok = l.peak()
	}
	return string(runes)
}

func (l *lexer) consumeNumber(first rune) string {
	integer := l.consumeDigits()

	// does the number contain a dot?
	r, ok := l.peak()
	if ok && r == '.' {
		// dot is part of the number only
		// when the dot is followed by a number.
		r, ok := l.peakTwice()
		if ok && isDigit(r) {
			// consume the dot
			_, _ = l.read()
			return string(first) + integer + "." + l.consumeDigits()
		}
	}
	return string(first) + integer
}

func (l *lexer) consumeIdentifier(first rune) string {
	runes := []rune{first}
	r, ok := l.peak()
	for ok && isAlphaNumeric(r) {
		r, _ = l.read()
		runes = append(runes, r)
		r, ok = l.peak()
	}
	return string(runes)
}

func (l *lexer) consumeString() (string, bool) {
	var runes []rune
	r, ok := l.peak()
	for ok && r != '"' {
		r, ok = l.read()
		if ok && r == '\n' {
			l.line++
		}
		runes = append(runes, r)
		r, ok = l.peak()
	}
	if r == '"' {
		_, _ = l.read()
		return string(runes), true
	}
	return string(runes), false
}

func (l *lexer) consumeWhiteSpace() string {
	var runes []rune
	r, ok := l.peak()
	for ok && (r == ' ' || r == '\t' || r == '\r') {
		r, _ = l.read()
		runes = append(runes, r)
		r, ok = l.peak()
	}
	return string(runes)
}

func (l *lexer) consumeLineComment() string {
	var runes []rune
	r, ok := l.peak()
	for ok && r != '\n' {
		r, ok = l.read()
		runes = append(runes, r)
		r, ok = l.peak()
	}
	return string(runes)
}

// Read consumes the next rune.
// Returns 0, false when the reader is at the end of the
// content or the next rune.
func (l *lexer) read() (rune, bool) {
	if l.readAheadBuffer.Len() > 0 {
		next := l.readAheadBuffer.Front()
		l.readAheadBuffer.Remove(next)
		r := next.Value.(rune)
		l.current += utf8.RuneLen(r)
		return r, true
	}
	r, n, err := l.content.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0, false
		}
		panic("unexpected error in read: " + err.Error())
	}
	l.current += n
	return r, true
}

func (l *lexer) peak() (rune, bool) {
	if l.readAheadBuffer.Len() == 0 {
		r, _, err := l.content.ReadRune()
		if err != nil {
			if err == io.EOF {
				return 0, false
			}
			panic("unexpected error in peak: " + err.Error())
		}
		_ = l.readAheadBuffer.PushBack(r)
		return r, true
	}
	r := l.readAheadBuffer.Front().Value.(rune)
	return r, true
}

func (l *lexer) peakTwice() (rune, bool) {
	if l.readAheadBuffer.Len() >= 2 {
		r := l.readAheadBuffer.Front().Next().Value.(rune)
		return r, true
	}
	if l.readAheadBuffer.Len() == 0 {
		r, _, err := l.content.ReadRune()
		if err != nil {
			if err == io.EOF {
				return 0, false
			}
			panic("unexpected error in peakTwice1: " + err.Error())
		}
		_ = l.readAheadBuffer.PushBack(r)
	}
	r, _, err := l.content.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0, false
		}
		panic("unexpected error in peakTwice2: " + err.Error())
	}
	_ = l.readAheadBuffer.PushBack(r)
	return r, true
}

func (l *lexer) isAtEnd() bool {
	return l.current >= l.size
}

func (l *lexer) next() Token {
	if l.isAtEnd() {
		return Token{
			TokenType: EOF,
			line:      l.line,
			lexeme:    "",
		}
	}
	return Token{
		TokenType: EOF,
		line:      l.line,
		lexeme:    "",
	}
}

func Lex(content string) ([]Token, error) {
	var res []Token
	l := lex(content)
	for {
		t := l.scanToken()
		switch t.TokenType {
		case ILLEGAL:
			return nil, errors.New("wrong")
		case EOF:
			return res, nil
		default:
			res = append(res, t)
		}
	}
}

func Parse(content string) error {
	return nil
}

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
