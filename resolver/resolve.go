package resolver

import (
	"fmt"

	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/types"
)

var handler *error.DiagnosticBag

type Table struct {
	symbols map[string]*Object
}

func InitTable() Table {
	t := Table{symbols: make(map[string]*Object, 4)}
	t.symbols["i8"] = newObj(TYPE)
	t.symbols["i8"].Type = types.Type{Alignment: 1, Size: 1, Base: nil}
	t.symbols["i16"] = newObj(TYPE)
	t.symbols["i16"].Type = types.Type{Alignment: 2, Size: 2, Base: nil}

	return t
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
	t.symbols[ast.Name] = newObj(VAR)
}
func (t *Table) isVariableExist(ident *ast.ExprIdent) {
	if _, ok := t.symbols[ident.Name]; !ok {
		handler.ReportError(error.Error{Msg: fmt.Sprintf("Variable '%s' doesn't exist", ident.Name), Pos: ident.Pos})
	}
}
func (t *Table) isTypeExist(typ types.TypeSpec) {
	if typ == nil {
		return
	}
	switch ty := typ.(type) {
	case *types.TypeName:
		{
			val, ok := t.symbols[ty.Name]
			if !ok {
				handler.ReportError(error.Error{Msg: fmt.Sprintf("Type '%s' doesn't exist", ty.Name)})
				return
			}
			if val.Kind != TYPE {
				handler.ReportError(error.Error{Msg: fmt.Sprintf("'%s' must be a type", ty.Name)})
			}

		}
	}
}

func Resolve(ast []ast.Decl, bag *error.DiagnosticBag) *Table {
	handler = bag
	table := InitTable()
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
			table.isTypeExist(n.RetType)
			for _, stmt := range n.Body {
				resolver(stmt, table)
			}
		}
	case *ast.StmtLet:
		{
			table.declareVariable(n)
			table.isTypeExist(n.Type)
			if n.Init != nil {
				resolver(n.Init, table)
			}

		}
	case *ast.StmtExpr:
		{
			resolver(n.Expr, table)

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
	case *ast.ExprInt:
		{

		}
	case *ast.ExprIdent:
		{
			table.isVariableExist(n)
		}

	}
}
