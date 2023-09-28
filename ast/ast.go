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
	Name     string
	TypeName string
	Init     Expr
}
type ExprBinary struct {
	Left  Expr
	Right Expr
	Op    byte // [0:'+',1:'*']
}
type ExprInt struct {
	Value int
}

func (e *ExprInt) exprNode()    {}
func (e *ExprBinary) exprNode() {}
func (s *StmtLet) stmtNode()    {}
