package resolver

import (
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/scope"
	"github.com/s0h1s2/types"
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
	Body       []StmtNode
}

type StmtLet struct {
	Name  string
	Pos   error.Position
	Init  ExprNode
	Scope *scope.Scope
	Type  *types.Type
}
type ExprLit struct {
	Value string
}
type ExprIdentifier struct {
	Name     string
	typee    *types.Type
	StackPos int
	Scope    *scope.Scope
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
