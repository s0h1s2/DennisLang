package parser

import "github.com/s0h1s2/lexer"
import "github.com/s0h1s2/error"

type Parser struct {
	tokens     []lexer.Token
	bag        *error.DiagnosticBag
	tokenIndex int
}

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
func (p *Parser) Parse() {
	println("----PARSER----")
	println(p.currentToken().String())

}
