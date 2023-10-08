package ast

import "github.com/s0h1s2/types"

type Node interface{}

type Expr interface {
	Node
	exprNode()
}
type Decl interface {
	Node
	declNode()
}
type DeclBad struct{}
type DeclFunction struct {
	Name    string
	RetType types.TypeSpec
	Body    []Stmt
}
type Stmt interface {
	Node
	stmtNode()
}

type StmtLet struct {
	Name string
	Type types.TypeSpec
	Init Expr
}

type ExprBinary struct {
	Left  Expr
	Right Expr
	Op    byte // [0:'+',1:'*']
}
type ExprAssign struct {
	Left  Expr
	Right Expr
}

type ExprIdent struct {
	Name string
}

type ExprInt struct {
	Value string
}

func (e *DeclFunction) declNode() {}
func (e *DeclBad) declNode()      {}
func (e *ExprInt) exprNode()      {}
func (e *ExprBinary) exprNode()   {}
func (e *ExprIdent) exprNode()    {}
func (e *ExprAssign) exprNode()   {}
func (s *StmtLet) stmtNode()      {}
