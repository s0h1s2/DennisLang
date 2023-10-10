package parser

import (
	"fmt"

	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/lexer"
	"github.com/s0h1s2/types"
)

type Parser struct {
	tokens     []lexer.Token
	bag        *error.DiagnosticBag
	tokenIndex int
	inRHS      bool
}
type Precedence byte

const (
	LOWEST Precedence = iota
	ASSIGN
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
	return &p.tokens[len(p.tokens)]
}
func (p *Parser) expectToken(kind lexer.TokenKind) *lexer.Token {
	if p.matchToken(kind) {
		token := p.currentToken()
		p.consumeToken()
		return token
	}
	p.reportHere(fmt.Sprintf("Expected '%s' but got '%s'", kind.String(), p.currentToken().Kind.String()))
	return nil
}
func (p *Parser) matchToken(kind lexer.TokenKind) bool {
	return kind == p.currentToken().Kind
}
func (p *Parser) consumeToken() {
	if p.tokenIndex < len(p.tokens) {
		p.tokenIndex++
	}
}
func (p *Parser) atEnd() bool {
	return p.currentToken().Kind == lexer.TK_EOF
}
func New(tokens []lexer.Token, bag *error.DiagnosticBag) *Parser {
	return &Parser{
		tokens:     tokens,
		bag:        bag,
		tokenIndex: 0,
		inRHS:      false,
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
	case lexer.TK_ASSIGN:
		return ASSIGN
	}
	return LOWEST

}
func (p *Parser) peekPreced() Precedence {
	prec := getPreced(p.peekToken())
	return prec
}
func (p *Parser) currPreced() Precedence {
	prec := getPreced(p.currentToken())
	if p.inRHS {
		return prec - 1
	}
	return prec
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
	p.reportHere(fmt.Sprintf("Unable to parse '%s'", p.currentToken().Kind.String()))
	return nil
}
func (p *Parser) parseBinary(left ast.Expr) ast.Expr {
	preced := p.currPreced()
	p.consumeToken()
	right := p.parseExpression(preced)
	return &ast.ExprBinary{Left: left, Right: right, Op: 1}
}
func (p *Parser) parseAssignment(left ast.Expr) ast.Expr {
	old := p.inRHS
	p.inRHS = true
	preced := p.currPreced()
	p.consumeToken()
	right := p.parseExpression(preced)
	p.inRHS = old
	return &ast.ExprAssign{Left: left, Right: right}
}

func (p *Parser) parseInfix(left ast.Expr) (ast.Expr, bool) {
	switch p.currentToken().Kind {
	case lexer.TK_PLUS:
		fallthrough
	case lexer.TK_STAR:
		{
			return p.parseBinary(left), true
		}
	case lexer.TK_ASSIGN:
		{
			return p.parseAssignment(left), true
		}
	}
	return nil, false
}
func (p *Parser) parseExpression(prec Precedence) ast.Expr {
	left := p.parseLeft()
	if left == nil {
		return nil
	}
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
func (p *Parser) parseVariableStmt() ast.Stmt {
	name := p.expectToken(lexer.TK_IDENT)
	p.expectToken(lexer.TK_COLON)
	typeSpec := p.parseType()
	var init ast.Expr
	if p.matchToken(lexer.TK_ASSIGN) {
		p.consumeToken()
		init = p.parseExpression(LOWEST)
		p.consumeToken()
	}
	p.expectToken(lexer.TK_SEMICOLON)
	return &ast.StmtLet{Name: name.Literal, Type: typeSpec, Init: init}
}
func (p *Parser) parseBlock() []ast.Stmt {
	p.expectToken(lexer.TK_OPENBRACE)
	stmts := []ast.Stmt{}
	for !p.atEnd() && p.currentToken().Kind != lexer.TK_CLOSEBRACE {
		switch p.currentToken().Kind {
		case lexer.TK_LET:
			{
				p.consumeToken()
				stmts = append(stmts, p.parseVariableStmt())
			}
		default:
			{
				p.reportHere(fmt.Sprintf("Unable to parse '%s' statement", p.currentToken().Kind.String()))
				p.consumeToken()
			}
		}
	}
	p.expectToken(lexer.TK_CLOSEBRACE)
	return stmts
}
func (p *Parser) parseBaseType() types.TypeSpec {
	name := p.expectToken(lexer.TK_IDENT)
	return &types.TypeName{Name: name.Literal}
}
func (p *Parser) parseType() types.TypeSpec {
	var left types.TypeSpec
	for p.matchToken(lexer.TK_STAR) {
		p.consumeToken()
		left = &types.TypePtr{Base: left}
	}
	if left != nil {
		switch t := left.(type) {
		case *types.TypePtr:
			{
				t.Base = p.parseBaseType()
			}
		}
	} else {
		left = p.parseBaseType()
	}

	return left
}

func (p *Parser) parseDeclarations() []ast.Decl {
	decls := []ast.Decl{}
	for p.currentToken().Kind != lexer.TK_EOF {
		switch p.currentToken().Kind {
		case lexer.TK_FN:
			{
				p.consumeToken()
				decls = append(decls, p.parseFunction())
			}
		default:
			{
				p.reportHere(fmt.Sprintf("Unable to parse '%s' declaration", p.currentToken().String()))
				p.consumeToken()
			}
		}
	}
	return decls
}
func (p *Parser) parseFunction() *ast.DeclFunction {
	name := p.expectToken(lexer.TK_IDENT)
	p.expectToken(lexer.TK_OPENPARAN)
	p.expectToken(lexer.TK_CLOSEPARAN)
	p.expectToken(lexer.TK_COLON)
	typeResult := p.parseType()
	body := p.parseBlock()
	return &ast.DeclFunction{Name: name.Literal, RetType: typeResult, Body: body}
}
func (p *Parser) Parse() []ast.Decl {
	println("----PARSER----")
	return p.parseDeclarations()
}
