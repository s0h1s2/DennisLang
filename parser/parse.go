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
	TERM              = iota
	FACTOR            = iota
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
	prev := p.currentToken()
	p.consumeToken()
	return &ast.ExprIdent{Name: prev.Literal}
}
func (p *Parser) parseInt() ast.Expr {
	prev := p.currentToken()
	p.consumeToken()
	return &ast.ExprInt{Value: prev.Literal}
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
	p.reportHere(fmt.Sprintf("Unknown kind of literal '%s'", p.currentToken().Kind.String()))
	return nil
}

func (p *Parser) parseExpression(prec Precedence) ast.Expr {
	left := p.parseLeft()
	return left
}

// func (p *Parser) parseBase() *ast.Expr {}
func (p *Parser) Parse() {
	println("----PARSER----")
	p.parseExpression(LOWEST)

}
