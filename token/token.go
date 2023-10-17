package token

import "fmt"
import "github.com/s0h1s2/error"

type TokenKind int

const (
	TK_INTEGER TokenKind = iota
	TK_IDENT
	TK_ASSIGN
	TK_PLUS
	TK_STAR
	TK_AND
	TK_ILLEGAL
	TK_SEMICOLON
	TK_COLON
	TK_OPENPARAN
	TK_CLOSEPARAN
	TK_OPENBRACE
	TK_CLOSEBRACE
	TK_EQUAL
	// KEYWORDS
	keywords_begin
	TK_LET
	TK_FN
	TK_RETURN
	TK_IF
	TK_TRUE
	TK_FALSE
	keywords_end

	TK_EOF
)

var tokenKindString = [...]string{
	TK_ILLEGAL:    "illegal",
	TK_PLUS:       "+",
	TK_STAR:       "*",
	TK_ASSIGN:     "=",
	TK_SEMICOLON:  ";",
	TK_COLON:      ":",
	TK_OPENPARAN:  "(",
	TK_CLOSEPARAN: ")",
	TK_OPENBRACE:  "{",
	TK_CLOSEBRACE: "}",
	TK_AND:        "&",
	TK_EQUAL:      "==",
	TK_INTEGER:    "integer",
	TK_IDENT:      "identifier",
	TK_RETURN:     "return",
	TK_LET:        "let",
	TK_FN:         "fn",
	TK_TRUE:       "true",
	TK_FALSE:      "false",
	TK_EOF:        "EOF",
}

func (tk TokenKind) String() string {
	return tokenKindString[tk]
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
