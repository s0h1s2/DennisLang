package ast

import "github.com/s0h1s2/types"
import "github.com/s0h1s2/error"

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
	Pos     error.Position
	Name    string
	RetType types.TypeSpec
	Body    []Stmt
	End     error.Position
}
type Stmt interface {
	Node
	stmtNode()
}

type StmtLet struct {
	Pos  error.Position
	Name string
	Type types.TypeSpec
	Init Expr
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
	Pos  error.Position
}

type ExprInt struct {
	Value string
}

func (e *DeclFunction) declNode() {}
func (e *DeclBad) declNode()      {}
func (s *StmtLet) stmtNode()      {}
func (s *StmtReturn) stmtNode()   {}
func (s *StmtExpr) stmtNode()     {}
func (e *ExprInt) exprNode()      {}
func (e *ExprBinary) exprNode()   {}
func (e *ExprIdent) exprNode()    {}
func (e *ExprAssign) exprNode()   {}
