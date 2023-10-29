package resolver

import (
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/scope"
	"github.com/s0h1s2/types"
)

type BinaryOperator = int

const (
	ADD BinaryOperator = iota
	SUB
	MUL
	DIV
	AND
	OR
)

type Node interface {
	GetType() *types.Type
	GetPos() error.Position
	GetScope() *scope.Scope
}

type DeclNode interface {
	Node
	declNode()
}
type StmtNode interface {
	Node
	stmtNode()
}
type ExprNode interface {
	Node
	exprNode()
}
type DeclFunction struct {
	Name       string
	ReturnType *types.Type
	StackSize  int
	Scope      *scope.Scope
	Body       StmtNode // StmtBlock
}
type Field struct {
	Name string
	Type *types.Type
}
type DeclStruct struct {
	Name   string
	Pos    error.Position
	Scope  *scope.Scope
	Fields []Field // TODO: fields need pos
}

type StmtLet struct {
	Name  string
	Pos   error.Position
	Init  ExprNode
	Scope *scope.Scope
	Type  *types.Type
}
type StmtExpr struct {
	Expr  ExprNode
	Scope *scope.Scope
}

type StmtBlock struct {
	Scope *scope.Scope
	Body  []StmtNode
}
type StmtReturn struct {
	Scope  *scope.Scope
	Result ExprNode
}
type ExprAssign struct {
	Left  ExprNode
	Right ExprNode
}
type ExprField struct {
	Name string
	Type *types.Type
	Pos  error.Position
}

type ExprInt struct {
	Value string
}
type ExprBool struct {
	Value string
}

type ExprIdentifier struct {
	Name string
	Type *types.Type
}

func (d *DeclFunction) declNode() {}
func (d *DeclFunction) GetType() *types.Type {
	return d.ReturnType
}
func (d *DeclFunction) GetPos() error.Position {
	return error.Position{}
}
func (d *DeclFunction) GetScope() *scope.Scope {
	return d.Scope
}

func (d *DeclStruct) declNode() {}
func (d *DeclStruct) GetType() *types.Type {
	return nil
}
func (d *DeclStruct) GetPos() error.Position {
	return error.Position{}
}
func (d *DeclStruct) GetScope() *scope.Scope {
	return d.Scope
}

func (d *StmtLet) stmtNode() {}
func (d *StmtLet) GetType() *types.Type {
	return d.Type
}
func (s *StmtLet) GetPos() error.Position {
	return error.Position{}
}
func (s *StmtLet) GetScope() *scope.Scope {
	return s.Scope
}
func (d *StmtBlock) stmtNode() {}
func (d *StmtBlock) GetType() *types.Type {
	return nil
}
func (s *StmtBlock) GetPos() error.Position {
	return error.Position{}
}
func (s *StmtBlock) GetScope() *scope.Scope {
	return s.Scope
}

func (d *StmtReturn) stmtNode() {}
func (d *StmtReturn) GetType() *types.Type {
	return nil
}
func (s *StmtReturn) GetPos() error.Position {
	return error.Position{}
}
func (s *StmtReturn) GetScope() *scope.Scope {
	return s.Scope
}

func (d *StmtExpr) stmtNode() {}
func (d *StmtExpr) GetType() *types.Type {
	return nil
}
func (s *StmtExpr) GetPos() error.Position {
	return error.Position{}
}
func (s *StmtExpr) GetScope() *scope.Scope {
	return s.Scope
}

func (e *ExprAssign) exprNode() {}
func (e *ExprAssign) GetType() *types.Type {
	return nil
}
func (e *ExprAssign) GetPos() error.Position {
	return error.Position{}
}
func (e *ExprAssign) GetScope() *scope.Scope {
	return nil
}
func (e *ExprIdentifier) exprNode() {}
func (e *ExprIdentifier) GetType() *types.Type {
	return e.Type
}
func (e *ExprIdentifier) GetPos() error.Position {
	return error.Position{}
}
func (e *ExprIdentifier) GetScope() *scope.Scope {
	return nil
}

func (e *ExprInt) exprNode() {}
func (e *ExprInt) GetType() *types.Type {
	return nil
}
func (e *ExprInt) GetPos() error.Position {
	return error.Position{}
}
func (e *ExprInt) GetScope() *scope.Scope {
	return nil
}

func (e *ExprField) exprNode() {}
func (e *ExprField) GetType() *types.Type {
	return e.Type
}
func (e *ExprField) GetPos() error.Position {
	return error.Position{}
}
func (e *ExprField) GetScope() *scope.Scope {
	return nil
}

func (e *ExprBool) exprNode() {}
func (e *ExprBool) GetType() *types.Type {
	return nil
}
func (e *ExprBool) GetPos() error.Position {
	return error.Position{}
}
func (e *ExprBool) GetScope() *scope.Scope {
	return nil
}
