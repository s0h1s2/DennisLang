package ast

type Visitor interface {
	Visit(node Node) (w Visitor)
}

func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}
	switch node.(type) {
	case *StmtLet:
		{
		}
	case *DeclFunction:
		{
		}
	default:
		{
		}
	}
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}
