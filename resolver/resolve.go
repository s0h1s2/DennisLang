package resolver

import (
	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/types"
)

var handler *error.DiagnosticBag

type Table struct {
	symbols *Scope
}

var table *Table

func InitTable() *Table {
	t := Table{symbols: NewScope(nil)}
	t.symbols.Define("i8", newObj(TYPE, types.NewType("i8", types.TYPE_INT, 1, 1)))
	t.symbols.Define("i16", newObj(TYPE, types.NewType("i16", types.TYPE_INT, 1, 1)))
	t.symbols.Define("i32", newObj(TYPE, types.NewType("i32", types.TYPE_INT, 1, 1)))
	t.symbols.Define("i64", newObj(TYPE, types.NewType("i64", types.TYPE_INT, 1, 1)))
	t.symbols.Define("bool", newObj(TYPE, types.NewType("bool", types.TYPE_BOOL, 1, 1)))
	t.symbols.Define("void", newObj(TYPE, types.NewType("void", types.TYPE_VOID, 0, 0)))
	return &t
}
func (t *Table) declareFunction(ast *ast.DeclFunction) {
	if result := t.symbols.Lookup(ast.Name); result {
		pos := ast.GetPos()
		handler.ReportError(pos, "Can't redeclare funciton '%s' more than once", ast.Name)
	}
	obj := newObj(FN, nil)
	obj.Node = ast
	t.symbols.Define(ast.Name, obj)
}
func (t *Table) GetObj(name string) *Object {
	if t.symbols.Lookup(name) {
		return t.symbols.GetObj(name)
	}
	return nil
}
func (t *Table) isVariableExist(ident *ast.ExprIdent) {
	if !t.symbols.Lookup(ident.Name) {
		pos := ident.GetPos()
		handler.ReportError(pos, "Variable '%s' doesn't exist", ident.Name)
	}
}
func (t *Table) isTypeExist(typ types.TypeSpec) (*types.Type, bool) {
	if typ == nil {
		return nil, false
	}

	pos := typ.GetPos()
	switch ty := typ.(type) {
	case *types.TypeName:
		{
			val := t.symbols.GetObj(ty.Name)
			if val == nil {
				handler.ReportError(pos, "Type '%s' doesn't exist", ty.Name)
				return nil, true
			}
			if val.Kind != TYPE {
				handler.ReportError(pos, "'%s' must be a type", ty.Name)
			}
			return val.Type, true
		}
	case *types.TypePtr:
		{
			if base, ok := t.isTypeExist(ty.Base); ok {
				ptr := types.NewType("*"+base.TypeName, types.TYPE_PTR, 8, 8)
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
	table = InitTable()
	for _, decl := range ast {
		resolver(decl, table.symbols)
	}
	return table
}
func resolver(node ast.Node, scope *Scope) {
	switch n := node.(type) {
	case *ast.DeclFunction:
		{
			table.declareFunction(n)
			if typ, ok := table.isTypeExist(n.RetType); ok {
				scope.GetObj(n.Name).Type = typ
				scope.GetObj(n.Name).Node = n
			}
			localScope := NewScope(scope)
			for _, stmt := range n.Body {
				resolver(stmt, localScope)

			}
		}
	case *ast.StmtLet:
		{
			if !scope.Lookup(n.Name) {
				obj := newObj(VAR, nil)
				if typ, ok := table.isTypeExist(n.Type); ok {
					obj.Type = typ
					scope.Define(n.Name, obj)
				}
			} else {
				handler.ReportError(n.GetPos(), "Can't redeclare variable '%s' more than one", n.Name)
			}
			if n.Init != nil {
				resolver(n.Init, scope)
			}

		}
	case *ast.StmtReturn:
		{
			resolver(n.Result, scope)
		}
	case *ast.StmtExpr:
		{
			resolver(n.Expr, scope)

		}
	case *ast.ExprBinary:
		{
			resolver(n.Left, scope)
			resolver(n.Right, scope)
		}
	case *ast.ExprAssign:
		{
			resolver(n.Left, scope)
			resolver(n.Right, scope)
		}
	case *ast.ExprInt:
		{

		}
	case *ast.ExprBoolean:
		{
		}
	case *ast.ExprIdent:
		{
			if !scope.Lookup(n.Name) {
				handler.ReportError(n.GetPos(), "Variable '%s' not found", n.Name)
			}
		}

	}
}
