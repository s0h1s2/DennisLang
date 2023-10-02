package ast

type Node interface{}

type Expr interface {
	Node
	exprNode()
}
type Stmt interface {
	Node
	stmtNode()
}

type StmtLet struct {
	Name string
	Type string
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

func (e *ExprInt) exprNode()    {}
func (e *ExprBinary) exprNode() {}
func (e *ExprIdent) exprNode()  {}
func (e *ExprAssign) exprNode() {}
func (s *StmtLet) stmtNode()    {}
