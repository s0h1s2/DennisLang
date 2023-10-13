package lexer

import (
	"fmt"

	"github.com/s0h1s2/error"
	"github.com/s0h1s2/token"
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

type Keyword = map[string]token.TokenKind

var keywords Keyword = Keyword{
	"if":     token.TK_IF,
	"let":    token.TK_LET,
	"fn":     token.TK_FN,
	"return": token.TK_RETURN,
}

func isKeyword(word string) token.TokenKind {
	if value, ok := keywords[word]; ok {
		return value
	}
	return token.TK_IDENT
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
func (lex *Lexer) makeToken(kind token.TokenKind, literal string) token.Token {
	return token.Token{
		Kind:    kind,
		Literal: literal,
		Pos:     error.Position{Start: lex.current, End: lex.current, Line: lex.line},
	}
}
func (lex *Lexer) atEnd() bool {
	return lex.current >= len(lex.src)
}
func (lex *Lexer) skipWhitespace() {
	for !lex.atEnd() && (lex.ch == ' ' || lex.ch == '\t') {
		lex.next()
	}
}
func (lex *Lexer) updateLine() {
	for !lex.atEnd() && lex.ch == '\n' {
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
func (lex *Lexer) getToken() token.Token {
start:
	lex.start = lex.current
	lex.updateLine()
	lex.skipWhitespace()
	if lex.atEnd() {
		return token.Token{Kind: token.TK_EOF}
	}

	lex.ch = lex.src[lex.current]
	for {
		switch lex.ch {
		case '+':
			{
				lex.next()
				return lex.makeToken(token.TK_PLUS, "")
			}
		case '*':
			{
				lex.next()
				return lex.makeToken(token.TK_STAR, "")
			}
		case '=':
			{
				lex.next()
				return lex.makeToken(token.TK_ASSIGN, "")
			}
		case '(':
			{
				lex.next()
				return lex.makeToken(token.TK_OPENPARAN, "")
			}
		case ')':
			{
				lex.next()
				return lex.makeToken(token.TK_CLOSEPARAN, "")
			}
		case '{':
			{
				lex.next()
				return lex.makeToken(token.TK_OPENBRACE, "")
			}
		case '}':
			{
				lex.next()
				return lex.makeToken(token.TK_CLOSEBRACE, "")
			}

		case ':':
			{
				lex.next()
				return lex.makeToken(token.TK_COLON, "")
			}
		case ';':
			{
				lex.next()
				return lex.makeToken(token.TK_SEMICOLON, "")
			}

		default:
			{
				if lex.ch >= '0' && lex.ch <= '9' {
					val := lex.scanInt()
					return lex.makeToken(token.TK_INTEGER, string(val))
				} else if isAlpha(lex.ch) || lex.ch == '_' {
					result := lex.scanIdentOrKeyword()
					return lex.makeToken(isKeyword(result), result)
				} else if lex.ch == '\n' {
					goto start
				}
				lex.errors.ReportError(error.Error{Msg: fmt.Sprintf("Illegal token '%c' with ascii code of '%d'", lex.src[lex.current], lex.src[lex.current]), Pos: error.Position{Line: lex.line, Start: lex.start, End: lex.current}})
				lex.next()
				return lex.makeToken(token.TK_ILLEGAL, "")
			}
		}
	}
}
func (lex *Lexer) GetTokens(src []byte) []token.Token {
	lex.reset()
	lex.src = src
	tk := lex.getToken()
	tokens := []token.Token{}
	for tk.Kind != token.TK_EOF {
		tokens = append(tokens, tk)
		tk = lex.getToken()
	}
	tokens = append(tokens, lex.makeToken(token.TK_EOF, ""))
	return tokens
}
