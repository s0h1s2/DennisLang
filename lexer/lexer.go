package lexer

import (
	"fmt"

	"github.com/s0h1s2/error"
)

type Lexer struct {
	src     []byte
	start   int
	current int
	line    int
	errors  *error.DiagnosticBag
}

func New(bag *error.DiagnosticBag) *Lexer {
	return &Lexer{
		src:     make([]byte, 0),
		start:   0,
		current: 0,
		line:    1,
		errors:  bag,
	}
}
func (lex *Lexer) reset() {
	lex.src = nil
	lex.start = 0
	lex.current = 0
}
func (lex *Lexer) next() {
	if lex.current < len(lex.src) {
		lex.current += 1
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
	return "0"
}
func (lex *Lexer) getToken() Token {
	lex.start = lex.current
	// TODO: jesus christ this is really bad!
	lex.skipWhitespace()
	lex.updateLine()
	lex.skipWhitespace()
	if lex.atEnd() {
		return Token{kind: TK_EOF, literal: "Hello"}
	}
	for {
		switch lex.src[lex.current] {
		case '+':
			{
				lex.next()
				return Token{kind: TK_PLUS, literal: "+"}
			}
		default:
			{
				lex.errors.ReportError(error.Error{Msg: fmt.Sprintf("Illegal token '%c'", lex.src[lex.current]), Pos: error.Position{Line: lex.line, Start: lex.start, End: lex.current}})
				lex.next()
				return Token{kind: TK_ILLEGAL, literal: "Illegal"}
			}
		}
	}
}
func (lex *Lexer) GetTokens(src []byte) []Token {
	lex.reset()
	lex.src = src
	token := lex.getToken()
	for token.kind != TK_EOF {
		println(token.kind)
		println(token.literal)
		token = lex.getToken()
	}
	return []Token{token}
}
