package lexer

import (
	"fmt"

	"github.com/s0h1s2/error"
)

type Lexer struct {
	src     []byte
	start   int
	current int
	ch      byte
	line    int
	errors  *error.DiagnosticBag
}

const (
	EOF = 0
)

type Keyword = map[string]TokenKind

var keywords Keyword = Keyword{
	"if":  TK_IF,
	"let": TK_LET,
}

func isKeyword(word string) TokenKind {
	if value, ok := keywords[word]; ok {
		return value
	}
	return TK_IDENT
}
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
func isNumeric(c byte) bool {
	return (c >= '0' && c <= '9')
}
func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isNumeric(c)
}
func New(bag *error.DiagnosticBag) *Lexer {
	return &Lexer{
		src:     make([]byte, 0),
		start:   0,
		current: 0,
		line:    1,
		ch:      ' ',
		errors:  bag,
	}
}
func (lex *Lexer) reset() {
	lex.src = nil
	lex.start = 0
	lex.current = 0
	lex.ch = 0
}
func (lex *Lexer) next() {
	if lex.current < len(lex.src) {
		lex.current += 1
	}
	lex.ch = EOF
	if !lex.atEnd() {
		lex.ch = lex.src[lex.current]
	}

}
func (lex *Lexer) makeToken(kind TokenKind, literal string) Token {
	return Token{
		Kind:    kind,
		Literal: literal,
		Pos:     error.Position{Start: lex.current, End: lex.current, Line: lex.line},
	}
}
func (lex *Lexer) atEnd() bool {
	return lex.current >= len(lex.src)
}
func (lex *Lexer) skipWhitespace() {
	for !lex.atEnd() && (lex.src[lex.current] == ' ' || lex.src[lex.current] == '\t') {
		lex.next()
	}
}
func (lex *Lexer) updateLine() {
	if !lex.atEnd() && lex.src[lex.current] == '\n' {
		lex.line += 1
		lex.next()
	}
}
func (lex *Lexer) scanInt() string {
	var val []byte
	for lex.ch >= '0' && lex.ch <= '9' {
		val = append(val, lex.ch)
		lex.next()
	}
	return string(val)
}
func (lex *Lexer) scanIdentOrKeyword() string {
	result := ""
	for isAlphaNumeric(lex.ch) || lex.ch == '_' {
		result += string(lex.ch)
		lex.next()
	}
	return result
}
func (lex *Lexer) getToken() Token {
	lex.start = lex.current
	// TODO: jesus christ this is really bad!
	lex.skipWhitespace()
	lex.updateLine()
	lex.skipWhitespace()
	if lex.atEnd() {
		return Token{Kind: TK_EOF}
	}
	lex.ch = lex.src[lex.current]
	for {
		switch lex.ch {
		case '+':
			{
				lex.next()
				return lex.makeToken(TK_PLUS, "")
			}
		case '*':
			{
				lex.next()
				return lex.makeToken(TK_STAR, "")
			}
		case '=':
			{
				lex.next()
				return lex.makeToken(TK_ASSIGN, "")
			}
		case ':':
			{
				lex.next()
				return lex.makeToken(TK_COLON, "")
			}
		case ';':
			{
				lex.next()
				return lex.makeToken(TK_SEMICOLON, "")
			}
		default:
			{
				if lex.ch >= '0' && lex.ch <= '9' {
					val := lex.scanInt()
					return lex.makeToken(TK_INTEGER, string(val))
				} else if isAlpha(lex.ch) || lex.ch == '_' {
					result := lex.scanIdentOrKeyword()
					return lex.makeToken(isKeyword(result), result)
				}
				lex.errors.ReportError(error.Error{Msg: fmt.Sprintf("Illegal token '%c'", lex.src[lex.current]), Pos: error.Position{Line: lex.line, Start: lex.start, End: lex.current}})
				lex.next()
				return lex.makeToken(TK_ILLEGAL, "")
			}
		}
	}
}
func (lex *Lexer) GetTokens(src []byte) []Token {
	lex.reset()
	lex.src = src
	token := lex.getToken()
	tokens := []Token{}
	for token.Kind != TK_EOF {
		tokens = append(tokens, token)
		token = lex.getToken()
	}
	return tokens
}
