package parser

import (
	"fmt"

	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/token"
	"github.com/s0h1s2/types"
)

type Parser struct {
	tokens     []token.Token
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

func (p *Parser) peekToken() *token.Token {
	if p.tokenIndex+1 < len(p.tokens) {
		return &p.tokens[p.tokenIndex+1]
	}
	return &p.tokens[p.tokenIndex]
}
func (p *Parser) currentToken() *token.Token {
	if p.tokenIndex < len(p.tokens) {
		return &p.tokens[p.tokenIndex]
	}
	return &p.tokens[len(p.tokens)]
}
func (p *Parser) expectToken(kind token.TokenKind) *token.Token {
	if p.matchToken(kind) {
		token := p.currentToken()
		p.consumeToken()
		return token
	}
	p.reportHere("Expected '%s' but got '%s'", kind.String(), p.currentToken().Kind.String())
	return nil
}
func (p *Parser) matchToken(kind token.TokenKind) bool {
	return kind == p.currentToken().Kind
}
func (p *Parser) consumeToken() {
	if p.tokenIndex < len(p.tokens) {
		p.tokenIndex++
	}
}
func (p *Parser) atEnd() bool {
	return p.currentToken().Kind == token.TK_EOF
}
func New(tokens []token.Token, bag *error.DiagnosticBag) *Parser {
	return &Parser{
		tokens:     tokens,
		bag:        bag,
		tokenIndex: 0,
		inRHS:      false,
	}
}
func (p *Parser) reportHere(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	p.bag.ReportError(error.Error{Msg: msg, Pos: p.currentToken().Pos})
}
func (p *Parser) parseIdent() ast.Expr {
	return &ast.ExprIdent{Name: p.currentToken().Literal, Pos: p.currentToken().Pos}
}
func (p *Parser) parseInt() ast.Expr {
	return &ast.ExprInt{Value: p.currentToken().Literal}
}
func getPreced(tk *token.Token) Precedence {
	switch tk.Kind {
	case token.TK_PLUS:
		return TERM
	case token.TK_STAR:
		return FACTOR
	case token.TK_ASSIGN:
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
func (p *Parser) parseBoolean() ast.Expr {
	val := false
	if p.currentToken().Kind == token.TK_TRUE {
		val = true
	}
	return &ast.ExprBoolean{Value: val}
}
func (p *Parser) parseLeft() ast.Expr {
	switch p.currentToken().Kind {
	case token.TK_IDENT:
		{
			return p.parseIdent()
		}
	case token.TK_INTEGER:
		{
			return p.parseInt()
		}
	case token.TK_AND:
		{
			return p.parseAddrOf()
		}
	case token.TK_FALSE:
		fallthrough
	case token.TK_TRUE:
		{
			return p.parseBoolean()
		}

	}
	return nil
}
func (p *Parser) parseAddrOf() ast.Expr {
	prevToken := p.currentToken()
	p.consumeToken()
	right := p.parseExpression(LOWEST)
	return &ast.ExprAddrOf{
		Pos:   prevToken.Pos,
		Right: right,
	}
}
func (p *Parser) parseBinary(left ast.Expr) ast.Expr {
	preced := p.currPreced()
	currentToken := p.currentToken()
	p.consumeToken()
	right := p.parseExpression(preced)
	return &ast.ExprBinary{Left: left, Right: right, Op: 1, Pos: currentToken.Pos}
}
func (p *Parser) parseAssignment(left ast.Expr) ast.Expr {
	old := p.inRHS
	p.inRHS = true
	currentToken := p.currentToken()
	preced := p.currPreced()
	p.consumeToken()
	right := p.parseExpression(preced)
	p.inRHS = old
	return &ast.ExprAssign{Left: left, Right: right, Pos: currentToken.Pos}
}

func (p *Parser) parseInfix(left ast.Expr) (ast.Expr, bool) {
	switch p.currentToken().Kind {
	case token.TK_PLUS:
		fallthrough
	case token.TK_STAR:
		{
			return p.parseBinary(left), true
		}
	case token.TK_ASSIGN:
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
func (p *Parser) parseReturn() ast.Stmt {
	prevToken := p.currentToken()
	result := p.parseExpression(LOWEST)
	if result != nil {
		p.consumeToken()
	}
	p.expectToken(token.TK_SEMICOLON)
	return &ast.StmtReturn{Result: result, Pos: prevToken.Pos}
}
func (p *Parser) parseVariableStmt() ast.Stmt {
	name := p.expectToken(token.TK_IDENT)
	p.expectToken(token.TK_COLON)
	typeSpec := p.parseType()
	var init ast.Expr
	if p.matchToken(token.TK_ASSIGN) {
		p.consumeToken()
		init = p.parseExpression(LOWEST)
		p.consumeToken()
	}
	p.expectToken(token.TK_SEMICOLON)
	return &ast.StmtLet{Name: name.Literal, Type: typeSpec, Init: init, Pos: name.Pos}
}
func (p *Parser) parseBlock() []ast.Stmt {
	p.expectToken(token.TK_OPENBRACE)
	stmts := []ast.Stmt{}
	for !p.atEnd() && p.currentToken().Kind != token.TK_CLOSEBRACE {
		switch p.currentToken().Kind {
		case token.TK_LET:
			{
				p.consumeToken()
				stmts = append(stmts, p.parseVariableStmt())
			}
		case token.TK_RETURN:
			{
				p.consumeToken()
				stmts = append(stmts, p.parseReturn())
			}
		default:
			{
				stmts = append(stmts, &ast.StmtExpr{Expr: p.parseExpression(LOWEST)})
				p.consumeToken()
				p.expectToken(token.TK_SEMICOLON)
			}
		}
	}
	p.expectToken(token.TK_CLOSEBRACE)
	return stmts
}
func (p *Parser) parseBaseType() types.TypeSpec {
	name := p.expectToken(token.TK_IDENT)
	if name == nil {
		return nil
	}
	return &types.TypeName{Name: name.Literal, Pos: name.Pos}
}
func (p *Parser) parseType() types.TypeSpec {
	var left types.TypeSpec
	prevToken := p.currentToken()
	for p.matchToken(token.TK_STAR) {
		p.consumeToken()
		left = &types.TypePtr{Base: left, Pos: prevToken.Pos}
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
	for p.currentToken().Kind != token.TK_EOF {
		switch p.currentToken().Kind {
		case token.TK_FN:
			{
				p.consumeToken()
				decls = append(decls, p.parseFunction())
			}
		default:
			{
				p.reportHere("Unable to parse '%s' declaration", p.currentToken().String())
				p.consumeToken()
			}
		}
	}
	return decls
}
func (p *Parser) parseFunction() *ast.DeclFunction {
	name := p.expectToken(token.TK_IDENT)
	p.expectToken(token.TK_OPENPARAN)
	p.expectToken(token.TK_CLOSEPARAN)
	p.expectToken(token.TK_COLON)
	typeResult := p.parseType()
	body := p.parseBlock()
	return &ast.DeclFunction{Name: name.Literal, RetType: typeResult, Body: body, Pos: name.Pos, End: p.currentToken().Pos}
}
func (p *Parser) Parse() []ast.Decl {
	println("----PARSER----")
	return p.parseDeclarations()
}
