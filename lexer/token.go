package lexer

import "fmt"
import "github.com/s0h1s2/error"

type TokenKind int

const (
	TK_INTEGER TokenKind = iota
	TK_IDENT
	TK_ASSIGN
	TK_PLUS
	TK_STAR
	TK_ILLEGAL
	TK_SEMICOLON
	TK_COLON
	// KEYWORDS
	keywords_begin
	TK_LET
	TK_IF
	keywords_end
	TK_EOF
)

func (tk TokenKind) String() string {
	switch tk {
	case TK_INTEGER:
		return "Integer"
	case TK_ILLEGAL:
		return "Illegal"
	case TK_PLUS:
		return "+"
	case TK_STAR:
		return "*"
	case TK_ASSIGN:
		return "="
	case TK_IDENT:
		return "Identifier"
	case TK_LET:
		return "let"
	case TK_SEMICOLON:
		return ";"
	case TK_COLON:
		return ":"
	case TK_EOF:
		return "EOF"

	}
	panic("Unreachable or unimplemented one of TokenKind")
}

type Token struct {
	Kind    TokenKind
	Literal string
	Pos     error.Position
}

func (tk *Token) String() string {
	lit := "nil"
	if tk.Literal != "" {
		lit = tk.Literal
	}
	return fmt.Sprintf("(%s , '%s')", tk.Kind.String(), lit)
}
