package resolver

import (
	"fmt"
	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
)

var handler *error.DiagnosticBag

type Table struct {
	symbols map[string]*Object
}

func (t *Table) declareFunction(ast *ast.DeclFunction) {
	if _, ok := t.symbols[ast.Name]; ok {
		handler.ReportError(error.Error{Msg: fmt.Sprintf("Can't redeclare funciton '%s' more than once", ast.Name), Pos: ast.Pos})
	}
	t.symbols[ast.Name] = newObj(FN)
}
func (t *Table) declareVariable(ast *ast.StmtLet) {
	if _, ok := t.symbols[ast.Name]; ok {
		handler.ReportError(error.Error{Msg: fmt.Sprintf("Can't redeclare variable '%s' more than once", ast.Name), Pos: ast.Pos})
	}
	t.symbols[ast.Name] = newObj(FN)
}
func (t *Table) isVariableExist(ident *ast.ExprIdent) {
	if _, ok := t.symbols[ident.Name]; !ok {
		handler.ReportError(error.Error{Msg: fmt.Sprintf("Variable '%s' doesn't exist", ident.Name), Pos: ident.Pos})
	}
}
func Resolve(ast []ast.Decl, bag *error.DiagnosticBag) *Table {
	handler = bag
	table := Table{
		symbols: make(map[string]*Object),
	}
	for _, decl := range ast {
		resolver(decl, &table)
	}
	return &table
}
func resolver(node ast.Node, table *Table) {
	switch n := node.(type) {
	case *ast.DeclFunction:
		{
			table.declareFunction(n)
			for _, stmt := range n.Body {
				resolver(stmt, table)
			}
		}
	case *ast.StmtLet:
		{
			table.declareVariable(n)
			resolver(n.Init, table)
		}
	case *ast.ExprBinary:
		{
			resolver(n.Left, table)
			resolver(n.Right, table)
		}
	case *ast.ExprAssign:
		{
			resolver(n.Left, table)
			resolver(n.Right, table)
		}
	case *ast.ExprIdent:
		{
			table.isVariableExist(n)
		}

	}
}
