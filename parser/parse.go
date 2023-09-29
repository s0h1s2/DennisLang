package parser

import (
	"fmt"

	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/lexer"
)

type Parser struct {
	tokens     []lexer.Token
	bag        *error.DiagnosticBag
	tokenIndex int
}
type Precedence byte

const (
	LOWEST Precedence = iota
	TERM
	FACTOR
)

func (p *Parser) peekToken() *lexer.Token {
	if p.tokenIndex+1 < len(p.tokens) {
		return &p.tokens[p.tokenIndex+1]
	}
	return &p.tokens[p.tokenIndex]
}
func (p *Parser) currentToken() *lexer.Token {
	if p.tokenIndex < len(p.tokens) {
		return &p.tokens[p.tokenIndex]
	}
	return &p.tokens[p.tokenIndex]
}
func (p *Parser) matchToken(kind lexer.TokenKind) bool {
	if kind == p.currentToken().Kind {
		return true
	}
	return false
}
func (p *Parser) consumeToken() {
	p.tokenIndex++
}
func (p *Parser) atEnd() bool {
	return p.currentToken().Kind == lexer.TK_EOF
}
func New(tokens []lexer.Token, bag *error.DiagnosticBag) *Parser {
	return &Parser{
		tokens:     tokens,
		bag:        bag,
		tokenIndex: 0,
	}
}
func (p *Parser) reportHere(msg string) {
	p.bag.ReportError(error.Error{Msg: msg, Pos: p.currentToken().Pos})
}
func (p *Parser) parseIdent() ast.Expr {
	return &ast.ExprIdent{Name: p.currentToken().Literal}
}
func (p *Parser) parseInt() ast.Expr {
	return &ast.ExprInt{Value: p.currentToken().Literal}
}
func getPreced(token *lexer.Token) Precedence {
	switch token.Kind {
	case lexer.TK_PLUS:
		return TERM
	case lexer.TK_STAR:
		return FACTOR
	}
	return LOWEST

}
func (p *Parser) peekPreced() Precedence {
	return getPreced(p.peekToken())
}
func (p *Parser) currPreced() Precedence {
	return getPreced(p.currentToken())
}

func (p *Parser) parseLeft() ast.Expr {
	switch p.currentToken().Kind {
	case lexer.TK_IDENT:
		{
			return p.parseIdent()
		}
	case lexer.TK_INTEGER:
		{
			return p.parseInt()
		}
	}
	p.reportHere(fmt.Sprintf("Unable to parse '%s'.", p.currentToken().Kind.String()))
	return nil
}
func (p *Parser) parseBinary(left ast.Expr) ast.Expr {
	preced := p.currPreced()
	p.consumeToken()
	right := p.parseExpression(preced)
	return &ast.ExprBinary{Left: left, Right: right, Op: 1}
}

func (p *Parser) parseInfix(left ast.Expr) (ast.Expr, bool) {
	switch p.currentToken().Kind {
	case lexer.TK_PLUS:
		fallthrough

	case lexer.TK_STAR:
		{
			return p.parseBinary(left), true
		}
	}
	return nil, false
}
func (p *Parser) parseExpression(prec Precedence) ast.Expr {
	left := p.parseLeft()
	for !p.atEnd() && prec < p.peekPreced() {
		p.consumeToken()
		right, ok := p.parseInfix(left)
		if !ok {
			return left
		}
		left = right
	}
	return left
}

func (p *Parser) Parse() {
	println("----PARSER----")
	result := p.parseExpression(LOWEST)
	switch result.(type) {
	case *ast.ExprInt:
		{

			println("int")
		}
	case *ast.ExprBinary:
		{

			println("binary")
		}
	default:
		{
			println("something else")
		}
	}
}
