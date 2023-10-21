package ast

import (
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/token"
	"github.com/s0h1s2/types"
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

type BinaryOpKind byte

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

type ExprIdent struct {
	Name string
	Pos  error.Position
}
type ExprAddrOf struct {
	Pos   error.Position
	Right Expr
}

type ExprInt struct {
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
func (s *StmtLet) stmtNode() {}
func (s *StmtLet) GetPos() error.Position {
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
func (e *ExprInt) exprNode() {}
func (e *ExprInt) GetPos() error.Position {
	return e.Pos
}

func (e *ExprBinary) exprNode() {}
func (e *ExprBinary) GetPos() error.Position {
	return e.Pos
}

func (e *ExprIdent) exprNode() {}
func (e *ExprIdent) GetPos() error.Position {
	return e.Pos
}

func (e *ExprAddrOf) exprNode() {}
func (e *ExprAddrOf) GetPos() error.Position {
	return e.Pos
}
func (e *ExprAssign) exprNode() {}
func (e *ExprAssign) GetPos() error.Position {
	return e.Pos
}

func (e *ExprBoolean) exprNode() {}
func (e *ExprBoolean) GetPos() error.Position {
	return e.Pos
}
