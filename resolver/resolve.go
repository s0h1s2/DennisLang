package resolver

import (
	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/scope"
	"github.com/s0h1s2/types"
)

var handler *error.DiagnosticBag

type Table struct {
	symbols *scope.Scope
}

var table *Table

func InitTable() *Table {
	t := Table{symbols: scope.NewScope(nil)}
	t.symbols.Define("i8", scope.NewTypeObj(types.NewType("i8", types.TYPE_INT, 1, 1)))
	t.symbols.Define("i16", scope.NewTypeObj(types.NewType("i16", types.TYPE_INT, 1, 1)))
	t.symbols.Define("i32", scope.NewTypeObj(types.NewType("i32", types.TYPE_INT, 1, 1)))
	t.symbols.Define("i64", scope.NewTypeObj(types.NewType("i64", types.TYPE_INT, 1, 1)))
	t.symbols.Define("bool", scope.NewTypeObj(types.NewType("bool", types.TYPE_BOOL, 1, 1)))
	t.symbols.Define("void", scope.NewTypeObj(types.NewType("void", types.TYPE_VOID, 0, 0)))
	return &t
}
func declareStruct(decl *ast.DeclStruct) {
	if table.symbols.Lookup(decl.Name) {
		handler.ReportError(decl.GetPos(), "Can't redeclare struct '%s' more than once", decl.Name)
	}
	structScope := scope.NewScope(nil)
	typ := types.NewType(decl.Name, types.TYPE_TYPE, 0, 0)
	table.symbols.Define(decl.Name, scope.NewObj(scope.TYPE, typ, structScope))
}
func (t *Table) GetScope() *scope.Scope {
	return t.symbols
}
func (t *Table) declareFunction(ast *ast.DeclFunction) {
	if t.symbols.Lookup(ast.Name) {
		pos := ast.GetPos()
		handler.ReportError(pos, "Can't redeclare funciton '%s' more than once", ast.Name)
	}
	obj := scope.NewObj(scope.FN, nil, nil)
	// obj.Node = ast
	t.symbols.Define(ast.Name, obj)
}
func (t *Table) GetObj(name string) *scope.Object {
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
			if val.Kind != scope.TYPE {
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
		resolver(decl, table.GetScope())
	}
	return table
}
func resolver(node ast.Node, currScope *scope.Scope) bool {
	switch n := node.(type) {
	case *ast.DeclFunction:
		{
			table.declareFunction(n)
			if typ, ok := table.isTypeExist(n.RetType); ok {
				currScope.GetObj(n.Name).Type = typ
				// currScope.GetObj(n.Name).Node = n
			}
			localScope := scope.NewScope(currScope)
			for _, stmt := range n.Body.Block {
				resolver(stmt, localScope)
			}
			table.GetObj(n.Name).Scope = localScope
			return true
		}
	case *ast.DeclStruct:
		{
			declareStruct(n)
			structScope := table.symbols.GetObj(n.Name).GetScope()
			for _, field := range n.Fields {
				ok := structScope.Define(field.Name, scope.NewObj(scope.VAR, nil, nil))
				if !ok {
					handler.ReportError(field.Pos, "Can't redeclare '%s' field in '%s'", field.Name, n.Name)
				} else {
					if typ, ok := table.isTypeExist(field.Type); ok {
						structScope.GetObj(field.Name).Type = typ
					}
				}
			}
			return true
		}
	case *ast.StmtLet:
		{
			if !currScope.Lookup(n.Name) {
				obj := scope.NewObj(scope.VAR, nil, currScope)
				if typ, ok := table.isTypeExist(n.Type); ok {
					obj.Type = typ
					currScope.Define(n.Name, obj)
				}
			} else {
				handler.ReportError(n.GetPos(), "Can't redeclare variable '%s' more than once", n.Name)
			}
			if n.Init != nil {
				resolver(n.Init, currScope)
			}

		}
	case *ast.StmtIf:
		{
			resolver(n.Cond, currScope)
			resolver(n.Then, currScope)
		}
	case *ast.StmtBlock:
		{
			newScope := scope.NewScope(currScope)
			for _, stmt := range n.Block {
				resolver(stmt, newScope)
			}
			n.Scope = newScope
		}
	case *ast.StmtReturn:
		{
			resolver(n.Result, currScope)
		}
	case *ast.StmtExpr:
		{
			resolver(n.Expr, currScope)

		}
	case *ast.ExprBinary:
		{
			resolver(n.Left, currScope)
			resolver(n.Right, currScope)
		}
	case *ast.ExprGet:
		{
			if currScope.Lookup(n.Name) {
				typ := table.symbols.GetObj(currScope.GetObj(n.Name).Type.TypeName)
				resolver(n.Right, typ.GetScope())
			}
		}

	case *ast.ExprAssign:
		{
			resolver(n.Left, currScope)
			resolver(n.Right, currScope)
		}
	case *ast.ExprInt:
		{
		}
	case *ast.ExprBoolean:
		{
		}
	case *ast.ExprIdent:
		{

			if !currScope.Lookup(n.Name) {
				handler.ReportError(n.GetPos(), "Variable '%s' not found", n.Name)
				// return false;
			}
			return true
		}
	default:
		{
			println(n)
			panic("Unreachable")
		}
	}
	return false
}
