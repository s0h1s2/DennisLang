package resolver

import (
	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/scope"
	"github.com/s0h1s2/types"
)

type Table struct {
	symbols *scope.Scope
}

var table *Table
var handler *error.DiagnosticBag

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
func Resolve(program []ast.Decl, bag *error.DiagnosticBag) *Table {
	println("----RESOLVER----")
	table = InitTable()
	handler = bag
	for _, decl := range program {
		resolveDecl(decl)
	}
	return table
}
func isTypeExist(typee types.TypeSpec) (*types.Type, bool) {
	switch t := typee.(type) {
	case *types.TypeName:
		{
			if table.symbols.Lookup(t.Name) {
				obj := table.symbols.GetObj(t.Name)
				if obj.Kind != scope.TYPE {
					handler.ReportError(typee.GetPos(), "Type '%s' must be a type not variable name or function name", t.Name)
					return nil, false
				}
				return obj.Type, true
			}
			handler.ReportError(typee.GetPos(), "Type '%s' doesn't exist", t.Name)
		}
	case *types.TypePtr:
		{
			isTypeExist(t.Base)
		}
	}

	return nil, false
}
func resolveDecl(decl ast.Decl) DeclNode {
	switch node := decl.(type) {
	case *ast.DeclFunction:
		{
			if table.symbols.LookupOnce(node.Name) {
				handler.ReportError(node.Pos, "Can't redeclare function '%s' more than once", node.Name)
				return nil
			}
			table.symbols.Define(node.Name, scope.NewObj(scope.FN, nil))
			resolveStmt(node.Body, nil)
		}
	}
	return &DeclFunction{}
}

func resolveStmt(stmt ast.Stmt, currScope *scope.Scope) StmtNode {
	pos := stmt.GetPos()
	switch node := stmt.(type) {
	case *ast.StmtLet:
		{
			if !currScope.LookupOnce(node.Name) {
				typ, ok := isTypeExist(node.Type)
				if !ok {
					return nil
				}
				currScope.Define(node.Name, scope.NewObj(scope.VAR, nil))
				if node.Init != nil {
					resolveExpr(node.Init, currScope)
				}
				return &StmtLet{Name: node.Name, Scope: currScope, Type: typ, Pos: node.Pos}
			}
			handler.ReportError(pos, "Can't redeclare '%s' variable more than once in same block", node.Name)
		}
	case *ast.StmtBlock:
		{
			s := scope.NewScope(currScope)
			for _, stmt := range node.Block {
				resolveStmt(stmt, s)
			}
		}
	}
	return nil
}
func resolveExpr(expr ast.Expr, scope *scope.Scope) {
	pos := expr.GetPos()
	switch node := expr.(type) {
	case *ast.ExprBinary:
		{
			resolveExpr(node.Left, scope)
			resolveExpr(node.Right, scope)
		}
	case *ast.ExprIdent:
		{
			if !scope.Lookup(node.Name) {
				handler.ReportError(pos, "Variable '%s' not found", node.Name)
			}
		}
	}
}

/*func declareStruct(decl *ast.DeclStruct) {
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
			if n.Result != nil {
				resolver(n.Result, currScope)
			}
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
}*/
