package lexer

import "fmt"
import "github.com/s0h1s2/error"

type TokenKind int

const (
	TK_INTEGER TokenKind = iota
	TK_IDENT
	TK_KEYWORD
	TK_PLUS
	TK_ILLEGAL
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
	case TK_EOF:
		return "EOF"

	}
	panic("Unreachable or unimplemented one of TokenKind")
}

type Token struct {
	kind    TokenKind
	literal string
	Pos     error.Position
}

func (tk *Token) String() string {
	lit := "nil"
	if tk.literal != "" || tk.kind != TK_ILLEGAL {
		lit = tk.literal
	}
	return fmt.Sprintf("(%s,%s)", tk.kind.String(), lit)
}
