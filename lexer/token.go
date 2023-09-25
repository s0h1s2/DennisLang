package lexer

type TokenKind int

const (
	TK_INTEGER TokenKind = iota
	TK_IDENT
	TK_KEYWORD
	TK_PLUS
	TK_ILLEGAL
	TK_EOF
)

type Token struct {
	kind    TokenKind
	literal string
}
