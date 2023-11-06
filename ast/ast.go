package ast

import (
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/scope"
	"github.com/s0h1s2/token"
)

type Node interface {
	GetPos() error.Position
}

type Expr interface {
	Node
	exprNode()
}
type Decl interface {
	Node
	declNode()
}
type DeclFunction struct {
	Pos        error.Position
	Name       string
	Parameters []Field
	RetType    TypeSpec
	Body       *StmtBlock
	End        error.Position
}
type Field struct {
	Pos  error.Position
	Name string
	Type TypeSpec
}
type DeclStruct struct {
	Pos    error.Position
	Name   string
	Fields []*Field
}
type Stmt interface {
	Node
	stmtNode()
}
type StmtBlock struct {
	Pos   error.Position
	Block []Stmt
	Scope *scope.Scope
}
type StmtLet struct {
	Pos  error.Position
	Name string
	Type TypeSpec
	Init Expr
}
type StmtIf struct {
	Pos  error.Position
	Cond Expr
	Then *StmtBlock
}

type StmtReturn struct {
	Pos    error.Position
	Result Expr
}

type StmtExpr struct {
	Pos  error.Position
	Expr Expr
}

type ExprBinary struct {
	Pos   error.Position
	Left  Expr
	Right Expr
	Op    token.TokenKind
}

type ExprAssign struct {
	Pos   error.Position
	Left  Expr
	Right Expr
}

type ExprCall struct {
	Pos  error.Position
	Name string
	Args []Expr
}

type CompoundField struct {
	Name string
	Init Expr
	Pos  error.Position
}
type ExprCompound struct {
	Pos    error.Position
	Type   TypeSpec
	Fields []CompoundField
}

type ExprField struct {
	Pos  error.Position
	Name string
	Expr Expr
}
type ExprIdent struct {
	Name string
	Pos  error.Position
}
type ExprUnary struct {
	Pos   error.Position
	Op    token.TokenKind
	Right Expr
}

type ExprInt struct {
	Pos   error.Position
	Value string
}
type ExprString struct {
	Pos   error.Position
	Value string
}

type ExprBoolean struct {
	Pos   error.Position
	Value bool
}

func (e *DeclFunction) declNode() {}
func (e *DeclFunction) GetPos() error.Position {
	return e.Pos
}
func (e *DeclStruct) declNode() {}
func (e *DeclStruct) GetPos() error.Position {
	return e.Pos
}
func (s *StmtLet) stmtNode() {}
func (s *StmtLet) GetPos() error.Position {
	return s.Pos
}
func (s *StmtIf) stmtNode() {}
func (s *StmtIf) GetPos() error.Position {
	return s.Pos
}
func (s *StmtBlock) stmtNode() {}
func (s *StmtBlock) GetPos() error.Position {
	return s.Pos
}

func (s *StmtReturn) stmtNode() {}
func (s *StmtReturn) GetPos() error.Position {
	return s.Pos
}
func (s *StmtExpr) stmtNode() {}
func (s *StmtExpr) GetPos() error.Position {
	return s.Pos
}
func (e *ExprBinary) exprNode() {}
func (e *ExprBinary) GetPos() error.Position {
	return e.Pos
}
func (e ExprCompound) exprNode() {}
func (e ExprCompound) GetPos() error.Position {
	return e.Pos
}
func (e *ExprUnary) exprNode() {}
func (e *ExprUnary) GetPos() error.Position {
	return e.Pos
}

func (e *ExprAssign) exprNode() {}
func (e *ExprAssign) GetPos() error.Position {
	return e.Pos
}
func (e *ExprField) exprNode() {}
func (e *ExprField) GetPos() error.Position {
	return e.Pos
}

func (e *ExprCall) exprNode() {}
func (e *ExprCall) GetPos() error.Position {
	return e.Pos
}

func (e *ExprIdent) exprNode() {}
func (e *ExprIdent) GetPos() error.Position {
	return e.Pos
}
func (e *ExprInt) exprNode() {}
func (e *ExprInt) GetPos() error.Position {
	return e.Pos
}

func (e *ExprBoolean) exprNode() {}
func (e *ExprBoolean) GetPos() error.Position {
	return e.Pos
}
func (e *ExprString) exprNode() {}
func (e *ExprString) GetPos() error.Position {
	return e.Pos
}
