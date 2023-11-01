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
	TK_DOT
	TK_BANG
	TK_EQUAL
	TK_NOTEQUAL
	TK_LESSTHAN
	TK_GREATERTHAN
	TK_GREATEREQUAL
	TK_LESSEQUAL
	// KEYWORDS
	keywords_begin
	TK_LET
	TK_FN
	TK_STRUCT
	TK_IF
	TK_TRUE
	TK_FALSE
	TK_RETURN
	keywords_end

	TK_EOF
)

var tokenKindString = [...]string{
	TK_ILLEGAL:      "illegal",
	TK_PLUS:         "+",
	TK_STAR:         "*",
	TK_ASSIGN:       "=",
	TK_SEMICOLON:    ";",
	TK_COLON:        ":",
	TK_OPENPARAN:    "(",
	TK_CLOSEPARAN:   ")",
	TK_OPENBRACE:    "{",
	TK_CLOSEBRACE:   "}",
	TK_AND:          "&",
	TK_BANG:         "!",
	TK_DOT:          ".",
	TK_LESSTHAN:     "<",
	TK_GREATERTHAN:  ">",
	TK_GREATEREQUAL: ">=",
	TK_LESSEQUAL:    "<=",
	TK_NOTEQUAL:     "!=",
	TK_EQUAL:        "==",
	TK_INTEGER:      "integer",
	TK_IDENT:        "identifier",
	TK_RETURN:       "return",
	TK_STRUCT:       "struct",
	TK_LET:          "let",
	TK_FN:           "fn",
	TK_IF:           "if",
	TK_TRUE:         "true",
	TK_FALSE:        "false",
	TK_EOF:          "EOF",
}

func (tk TokenKind) String() string {
	return tokenKindString[tk]
}

type Token struct {
	Kind    TokenKind
	Literal string
	Pos     error.Position
}
type keywordMap = map[string]TokenKind

func InitKeywords() keywordMap {
	keywords := make(keywordMap, (keywords_end-keywords_begin)+1)
	for i := keywords_begin + 1; i < keywords_end; i++ {
		keywords[tokenKindString[i]] = i
	}
	return keywords
}
func (tk *Token) String() string {
	lit := "nil"
	if tk.Literal != "" {
		lit = tk.Literal
	}
	return fmt.Sprintf("(%s , '%s')", tk.Kind.String(), lit)
}
