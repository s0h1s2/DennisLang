package resolver

import (
	"fmt"

	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/scope"
	"github.com/s0h1s2/types"
)

type Table struct {
	Symbols *scope.Scope
}

var table *Table
var handler *error.DiagnosticBag
var isFieldAccess bool

func InitTable() *Table {
	t := Table{Symbols: scope.NewScope(nil)}
	t.Symbols.Define("i8", scope.NewTypeObj(types.NewType("i8", types.TYPE_INT, 1, 1)))
	t.Symbols.Define("i16", scope.NewTypeObj(types.NewType("i16", types.TYPE_INT, 1, 1)))
	t.Symbols.Define("i32", scope.NewTypeObj(types.NewType("i32", types.TYPE_INT, 1, 1)))
	t.Symbols.Define("i64", scope.NewTypeObj(types.NewType("i64", types.TYPE_INT, 1, 1)))
	t.Symbols.Define("bool", scope.NewTypeObj(types.NewType("bool", types.TYPE_BOOL, 1, 1)))
	t.Symbols.Define("void", scope.NewTypeObj(types.NewType("void", types.TYPE_VOID, 0, 0)))
	return &t
}
func Resolve(program []ast.Decl, bag *error.DiagnosticBag) (*Table, []DeclNode) {
	println("----RESOLVER----")
	table = InitTable()
	handler = bag
	var decls []DeclNode
	for _, decl := range program {
		decls = append(decls, resolveDecl(decl))
	}
	return table, decls
}
func isTypeExist(typee types.TypeSpec) (*types.Type, bool) {
	switch t := typee.(type) {
	case *types.TypeName:
		{
			if table.Symbols.Lookup(t.Name) {
				obj := table.Symbols.GetObj(t.Name)
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
			if table.Symbols.LookupOnce(node.Name) {
				handler.ReportError(node.Pos, "Can't redeclare function '%s' more than once", node.Name)
				return nil
			}
			typ, ok := isTypeExist(node.RetType)
			if !ok {
				return nil
			}
			table.Symbols.Define(node.Name, scope.NewObj(scope.FN, nil))
			resolvedBody := resolveStmt(node.Body, nil)
			return &DeclFunction{Scope: resolvedBody.GetScope(), Name: node.Name, Body: resolvedBody, ReturnType: typ}
		}
	case *ast.DeclStruct:
		{
			if table.Symbols.LookupOnce(node.Name) {
				handler.ReportError(node.Pos, "Can't redeclare struct '%s' more than once", node.Name)
				return nil
			}

			structScope := scope.NewScope(nil)
			obj := scope.NewObj(scope.TYPE, types.NewType(node.Name, types.TYPE_TYPE, 0, 0))
			obj.Scope = structScope
			table.Symbols.Define(node.Name, obj)
			fields := make([]Field, 0, 4)
			for _, field := range node.Fields {
				if structScope.LookupOnce(field.Name) {
					handler.ReportError(field.Pos, "Can't redeclare '%s' field more than once in struct '%s'", field.Name, node.Name)
					return nil
				}

				typ, ok := isTypeExist(field.Type)
				if !ok {
					return nil
				}
				obj := scope.NewObj(scope.FIELD, typ)
				if typ.Kind == types.TYPE_TYPE {
					obj.Scope = table.Symbols.GetObj(typ.TypeName).Scope
				}
				structScope.Define(field.Name, obj)
				fields = append(fields, Field{Name: field.Name, Type: typ})
			}
			return &DeclStruct{Name: node.Name, Fields: fields, Pos: node.Pos, Scope: structScope}
		}
	}
	return nil
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
				currScope.Define(node.Name, scope.NewObj(scope.VAR, typ))
				var resolvedExpr ExprNode
				if node.Init != nil {
					resolvedExpr = resolveExpr(node.Init, currScope, nil)
				}

				return &StmtLet{Name: node.Name, Init: resolvedExpr, Scope: currScope, Type: typ, Pos: node.Pos}
			}
			handler.ReportError(pos, "Can't redeclare '%s' variable more than once in same block", node.Name)
		}
	case *ast.StmtReturn:
		{
			if node.Result != nil {
				resolvedExpr := resolveExpr(node.Result, currScope, nil)
				return &StmtReturn{Result: resolvedExpr}
			}
		}
	case *ast.StmtExpr:
		{
			// TODO: return resolved stmt expr
			expr := resolveExpr(node.Expr, currScope, nil)
			return &StmtExpr{Expr: expr, Scope: currScope}
		}
	case *ast.StmtBlock:
		{
			s := scope.NewScope(currScope)
			var resolvedStmts []StmtNode
			for _, stmt := range node.Block {
				resolvedStmts = append(resolvedStmts, resolveStmt(stmt, s))
			}
			return &StmtBlock{Scope: s, Body: resolvedStmts}
		}
	}
	return nil
}
func resolveExpr(expr ast.Expr, scope *scope.Scope, typeScope *scope.Scope) ExprNode {
	pos := expr.GetPos()
	switch node := expr.(type) {
	case *ast.ExprBinary:
		{
			resolveExpr(node.Left, scope, nil)
			resolveExpr(node.Right, scope, nil)
		}
	case *ast.ExprAssign:
		{
			left := resolveExpr(node.Left, scope, nil)
			right := resolveExpr(node.Right, scope, nil)
			return &ExprAssign{Right: right, Left: left}
		}
	case *ast.ExprGet:
		{
			resolveExpr(node.Name, scope, nil)
			resolveExpr(node.Right, scope, nil)
		}

	case *ast.ExprInt:
		{
			return &ExprInt{Value: node.Value}
		}
	case *ast.ExprBoolean:
		{
			return &ExprBool{Value: "1"}
		}
	case *ast.ExprField:
		{
			println(node.Name)
		}
	case *ast.ExprIdent:
		{
			if !scope.Lookup(node.Name) {
				handler.ReportError(pos, "Variable '%s' not found", node.Name)
				return nil
			}
			return &ExprIdentifier{Name: node.Name, Type: scope.GetObj(node.Name).Type}
		}
	default:
		{
			panic(fmt.Sprintf("Unhandled node '%T' or Unreachable", node))
		}
	}
	return nil
}
