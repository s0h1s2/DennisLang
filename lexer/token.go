package lexer

import "fmt"

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
	case TK_EOF:
		return "EOF"

	}
	return "Unreachable"
}

type Token struct {
	kind    TokenKind
	literal string
}

func (tk Token) String() string {
	return fmt.Sprintf("(%s,%s)", tk.kind.String(), tk.literal)
}
