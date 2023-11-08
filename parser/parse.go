/*
May you do good and not evil.
May you find forgiveness for yourself and forgive others.
May you share freely, never taking more than you give.
*/
package parser

import (
	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/token"
)

type Parser struct {
	tokens        []token.Token
	bag           *error.DiagnosticBag
	tokenIndex    int
	inRHS         bool
	hadError      bool
	isFlowControl bool
}

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
	p.hadError = true
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
		tokens:        tokens,
		bag:           bag,
		tokenIndex:    0,
		hadError:      false,
		isFlowControl: false,
	}
}
func (p *Parser) reportHere(format string, args ...interface{}) {
	p.bag.ReportError(p.currentToken().Pos, format, args...)
}
func (p *Parser) parseIdent() ast.Expr {
	return &ast.ExprIdent{Name: p.currentToken().Literal, Pos: p.currentToken().Pos}
}
func (p *Parser) parseInt() ast.Expr {
	return &ast.ExprInt{Value: p.currentToken().Literal}
}
func (p *Parser) parseBoolean() ast.Expr {
	val := false
	if p.currentToken().Kind == token.TK_TRUE {
		val = true
	}
	return &ast.ExprBoolean{Value: val}
}
func (p *Parser) parseCompound(typ ast.TypeSpec) ast.Expr {
	tk := p.expectToken(token.TK_OPENBRACE)
	fields := make([]ast.CompoundField, 0, 4)
	for !p.matchToken(token.TK_CLOSEBRACE) {
		name := p.expectToken(token.TK_IDENT)
		if name == nil {
			p.reportHere("Expected field name in compound initialization")
			return nil
		}
		p.expectToken(token.TK_COLON)
		init := p.parseExpression()
		fields = append(fields, ast.CompoundField{Name: name.Literal, Init: init, Pos: name.Pos})
		if !p.matchToken(token.TK_COMMA) {
			break
		}
		p.consumeToken()
	}
	p.expectToken(token.TK_CLOSEBRACE)
	return &ast.ExprCompound{Type: typ, Fields: fields, Pos: tk.Pos}
}
func (p *Parser) parsePrimary() ast.Expr {
	switch p.currentToken().Kind {
	case token.TK_IDENT:
		{
			ident := p.parseIdent()
			p.consumeToken()
			if p.matchToken(token.TK_OPENBRACE) && !p.isFlowControl {
				typeName := ident.(*ast.ExprIdent)
				return p.parseCompound(&ast.TypeName{Name: typeName.Name, Pos: typeName.Pos})
			}
			return ident
		}
	case token.TK_INTEGER:
		{
			integer := p.parseInt()
			p.consumeToken()
			return integer
		}
	case token.TK_TRUE:
		fallthrough
	case token.TK_FALSE:
		{
			boolean := p.parseBoolean()
			p.consumeToken()
			return boolean
		}
	case token.TK_STRING:
		{
			tk := p.currentToken()
			str := &ast.ExprString{Value: tk.Literal, Pos: tk.Pos}
			p.consumeToken()
			return str
		}
	case token.TK_OPENPARAN:
		{
			p.consumeToken()
			expr := p.parseExpression()
			p.expectToken(token.TK_CLOSEPARAN)
			return expr
		}
	default:
		{
			p.reportHere("Unexpected token '%s' in expression", p.currentToken().Kind.String())
		}
	}
	p.consumeToken()
	return nil
}
func (p *Parser) parseBase() ast.Expr {
	expr := p.parsePrimary()
	for p.matchToken(token.TK_OPENPARAN) || p.matchToken(token.TK_DOT) {
		if p.matchToken(token.TK_DOT) {
			p.consumeToken()
			name := p.expectToken(token.TK_IDENT)
			expr = &ast.ExprField{Expr: expr, Name: name.Literal, Pos: name.Pos}
		} else if p.matchToken(token.TK_OPENPARAN) {
			paranPos := p.expectToken(token.TK_OPENPARAN)
			// parse function arguments
			args := make([]ast.Expr, 0)
			for !p.matchToken(token.TK_CLOSEPARAN) {
				args = append(args, p.parseExpression())
				if !p.matchToken(token.TK_COMMA) {
					break
				}
				p.consumeToken()
			}
			p.expectToken(token.TK_CLOSEPARAN)
			name, ok := expr.(*ast.ExprIdent)
			if !ok {
				p.reportHere("Function call must be a name")
				return nil
			}
			expr = &ast.ExprCall{Pos: paranPos.Pos, Args: args, Name: name.Name}
		}
	}
	return expr
}
func (p *Parser) parseUnary() ast.Expr {
	if p.matchToken(token.TK_AND) || p.matchToken(token.TK_STAR) || p.matchToken(token.TK_BANG) {
		op := p.currentToken()
		p.consumeToken()
		return &ast.ExprUnary{Op: op.Kind, Pos: op.Pos, Right: p.parseUnary()}
	} else {
		return p.parseBase()
	}
}
func (p *Parser) parseFactor() ast.Expr {
	left := p.parseUnary()
	for p.matchToken(token.TK_STAR) {
		op := p.currentToken()
		p.consumeToken()
		left = &ast.ExprBinary{Left: left, Right: p.parseUnary(), Op: op.Kind, Pos: op.Pos}
	}
	return left

}
func (p *Parser) parseTerm() ast.Expr {
	left := p.parseFactor()
	for p.matchToken(token.TK_PLUS) {
		op := p.currentToken()
		p.consumeToken()
		left = &ast.ExprBinary{Left: left, Right: p.parseFactor(), Op: op.Kind, Pos: op.Pos}
	}
	return left

}

func (p *Parser) parseCompare() ast.Expr {
	left := p.parseTerm()
	for p.matchToken(token.TK_EQUAL) || p.matchToken(token.TK_NOTEQUAL) || p.matchToken(token.TK_LESSTHAN) || p.matchToken(token.TK_LESSEQUAL) || p.matchToken(token.TK_GREATEREQUAL) || p.matchToken(token.TK_GREATERTHAN) {
		op := p.currentToken()
		p.consumeToken() // Consume operator
		left = &ast.ExprBinary{Right: p.parseTerm(), Op: op.Kind, Left: left, Pos: op.Pos}
	}
	return left
}
func (p *Parser) parseAssignment() ast.Expr {
	left := p.parseCompare()
	if p.matchToken(token.TK_ASSIGN) {
		assign := p.currentToken()
		p.consumeToken()
		left = &ast.ExprAssign{Left: left, Right: p.parseExpression(), Pos: assign.Pos}
	}
	return left
}
func (p *Parser) parseExpression() ast.Expr {
	return p.parseAssignment()
}
func (p *Parser) parseVariableStmt() ast.Stmt {
	name := p.expectToken(token.TK_IDENT)
	p.expectToken(token.TK_COLON)
	typeSpec := p.parseType()
	var init ast.Expr
	if p.matchToken(token.TK_ASSIGN) {
		p.consumeToken()
		init = p.parseExpression()
	}
	p.expectToken(token.TK_SEMICOLON)
	return &ast.StmtLet{Name: name.Literal, Type: typeSpec, Init: init, Pos: name.Pos}
}
func (p *Parser) parseReturn() ast.Stmt {
	ret := p.expectToken(token.TK_RETURN)
	var expr ast.Expr
	if !p.matchToken(token.TK_SEMICOLON) {
		expr = p.parseExpression()
	}
	p.expectToken(token.TK_SEMICOLON)
	return &ast.StmtReturn{Pos: ret.Pos, Result: expr}
}
func (p *Parser) parseLoop() ast.Stmt {
	loopToken := p.expectToken(token.TK_LOOP)
	p.isFlowControl = true
	expr := p.parseExpression()
	p.isFlowControl = false
	then := p.parseBlock()
	return &ast.StmtLoop{Cond: expr, Then: then, Pos: loopToken.Pos}
}
func (p *Parser) parseIf() ast.Stmt {
	p.expectToken(token.TK_IF)
	pos := p.currentToken().Pos
	p.isFlowControl = true
	cond := p.parseExpression()
	p.isFlowControl = false
	then := p.parseBlock()
	var elseBody *ast.StmtBlock = nil
	var elseIf ast.Stmt = nil
	if p.matchToken(token.TK_ELSE) {
		p.consumeToken()
		if p.matchToken(token.TK_IF) {
			elseIf = p.parseIf()
		} else {
			elseBody = p.parseBlock()
		}
	}
	return &ast.StmtIf{
		Cond:   cond,
		Then:   then,
		ElseIf: elseIf,
		Else:   elseBody,
		Pos:    pos,
	}
}
func (p *Parser) parseBlock() *ast.StmtBlock {
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
				stmts = append(stmts, p.parseReturn())
			}
		case token.TK_OPENBRACE:
			{
				stmts = append(stmts, p.parseBlock())
			}
		case token.TK_IF:
			{
				stmts = append(stmts, p.parseIf())
			}
		case token.TK_LOOP:
			{
				stmts = append(stmts, p.parseLoop())
			}
		default:
			{
				stmts = append(stmts, &ast.StmtExpr{Expr: p.parseExpression()})
				p.expectToken(token.TK_SEMICOLON)
			}
		}
	}

	p.expectToken(token.TK_CLOSEBRACE)
	return &ast.StmtBlock{Block: stmts}
}
func (p *Parser) parseBaseType() ast.TypeSpec {
	if p.matchToken(token.TK_IDENT) || p.matchToken(token.TK_STRING) {
		name := p.expectToken(p.currentToken().Kind)
		return &ast.TypeName{Name: name.Literal, Pos: name.Pos}
	}
	p.reportHere("Expected type but got '%s'", p.currentToken().Kind.String())
	return nil
}
func (p *Parser) parseType() ast.TypeSpec {
	prevToken := p.currentToken()
	if p.matchToken(token.TK_STAR) {
		p.consumeToken()
		left := &ast.TypePtr{Base: p.parseType(), Pos: prevToken.Pos}
		return left
	} else {
		return p.parseBaseType()
	}

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
		case token.TK_STRUCT:
			{
				p.consumeToken()
				decls = append(decls, p.parseStruct())
			}
		case token.TK_EXTERN:
			{
				p.consumeToken()
				decls = append(decls, p.parseExternal())
			}
		default:
			{
				p.reportHere("Unable to parse '%s' declaration", p.currentToken().Kind.String())
				p.consumeToken()
			}
		}
	}
	return decls
}
func (p *Parser) hasError() bool {
	if p.hadError {
		p.hadError = false
		return true
	}
	return false
}
func (p *Parser) parseField() *ast.Field {
	name := p.expectToken(token.TK_IDENT)
	p.expectToken(token.TK_COLON)
	typ := p.parseType()
	if p.hasError() {
		return nil
	}
	return &ast.Field{Name: name.Literal, Type: typ, Pos: name.Pos}
}
func (p *Parser) parseFunctionHeader() (*token.Token, []ast.Field, ast.TypeSpec) {
	name := p.expectToken(token.TK_IDENT)
	p.expectToken(token.TK_OPENPARAN)
	params := p.parseFunctionParameters()
	p.expectToken(token.TK_CLOSEPARAN)
	p.expectToken(token.TK_COLON)
	typeResult := p.parseType()
	return name, params, typeResult
}
func (p *Parser) parseExternal() ast.Decl {
	if p.matchToken(token.TK_FN) {
		p.consumeToken()
		name, params, typee := p.parseFunctionHeader()
		p.expectToken(token.TK_SEMICOLON)
		return &ast.DeclExternalFunction{Parameters: params, Name: name.Literal, ReturnType: typee, Pos: name.Pos}
	}
	p.reportHere("Expected 'fn' or 'let' after 'external' keyword but got '%s'", p.currentToken().Kind.String())
	return nil
}
func (p *Parser) parseStruct() ast.Decl {
	name := p.expectToken(token.TK_IDENT)
	p.expectToken(token.TK_OPENBRACE)
	fields := make([]*ast.Field, 0, 4)
	for !p.atEnd() && !p.matchToken(token.TK_CLOSEBRACE) {
		fields = append(fields, p.parseField())
		if p.expectToken(token.TK_SEMICOLON) == nil {
			return nil
		}
	}
	p.expectToken(token.TK_CLOSEBRACE)

	if p.hasError() {
		return nil
	}
	return &ast.DeclStruct{
		Name:   name.Literal,
		Fields: fields,
		Pos:    name.Pos,
	}
}

func (p *Parser) parseFunctionParameters() []ast.Field {
	params := make([]ast.Field, 0)
	for !p.atEnd() && !p.matchToken(token.TK_CLOSEPARAN) {
		name := p.expectToken(token.TK_IDENT)
		if name == nil {
			return nil
		}
		p.expectToken(token.TK_COLON)
		typeSpec := p.parseType()
		params = append(params, ast.Field{Name: name.Literal, Type: typeSpec})
		if !p.matchToken(token.TK_COMMA) {
			break
		}
		p.consumeToken()
	}
	return params
}
func (p *Parser) parseFunction() *ast.DeclFunction {
	name, params, typeResult := p.parseFunctionHeader()
	body := p.parseBlock()
	if p.hasError() {
		return nil
	}
	return &ast.DeclFunction{Name: name.Literal, RetType: typeResult, Body: body, Pos: name.Pos, Parameters: params, End: p.currentToken().Pos}
}
func (p *Parser) Parse() []ast.Decl {
	println("----PARSER----")
	return p.parseDeclarations()
}
