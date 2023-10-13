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
	t.symbols["i8"].Type = *types.NewType(1, 1)
	t.symbols["i16"] = newObj(TYPE)
	t.symbols["i16"].Type = *types.NewType(2, 2)
	return t
}

func (t *Table) declareFunction(ast *ast.DeclFunction) {
	if _, ok := t.symbols[ast.Name]; ok {
		handler.ReportError(error.Error{Msg: fmt.Sprintf("Can't redeclare funciton '%s' more than once", ast.Name), Pos: ast.Pos})
	}
	t.symbols[ast.Name] = newObj(FN)
}
func (t *Table) getObj(name string) *Object {
	return t.symbols[name]
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
func (t *Table) isTypeExist(typ types.TypeSpec) (*types.Type, bool) {
	if typ == nil {
		return nil, false
	}
	switch ty := typ.(type) {
	case *types.TypeName:
		{
			val, ok := t.symbols[ty.Name]
			if !ok {
				handler.ReportError(error.Error{Msg: fmt.Sprintf("Type '%s' doesn't exist", ty.Name)})
				return nil, true
			}
			if val.Kind != TYPE {
				handler.ReportError(error.Error{Msg: fmt.Sprintf("'%s' must be a type", ty.Name)})
			}
			return &val.Type, true
		}
	case *types.TypePtr:
		{
			if base, ok := t.isTypeExist(ty.Base); ok {
				ptr := types.NewType(8, 8)
				ptr.Base = base
				return ptr, true

			}
		}
	}
	return nil, false
}

func Resolve(ast []ast.Decl, bag *error.DiagnosticBag) *Table {
	println("----RESOLVER----")
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
			if typ, ok := table.isTypeExist(n.RetType); ok {
				table.getObj(n.Name).Type = *typ
				table.getObj(n.Name).Node = n
			}
			for _, stmt := range n.Body {
				resolver(stmt, table)
			}
		}
	case *ast.StmtLet:
		{
			table.declareVariable(n)
			if typ, ok := table.isTypeExist(n.Type); ok {
				table.getObj(n.Name).Type = *typ
				table.getObj(n.Name).Node = n
			}
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
