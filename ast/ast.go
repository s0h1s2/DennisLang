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

type ExprInt struct {
	Value int
}

func (e *ExprInt) exprNode() {}
func (s *StmtLet) stmtNode() {}
