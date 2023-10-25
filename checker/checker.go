package checker

/*package checker

import (
	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/resolver"
	"github.com/s0h1s2/scope"
	"github.com/s0h1s2/token"
	"github.com/s0h1s2/types"
)

var symTable *resolver.Table
var currScope *scope.Scope
var prevScope *scope.Scope
var declerations []ast.Decl
var handler *error.DiagnosticBag
var functionName string

func enterScope(scope *scope.Scope) {
	prevScope = currScope
	currScope = scope
}
func leaveScope() {
	currScope = prevScope
}

func TypeChecker(table *resolver.Table, decls []ast.Decl, bag *error.DiagnosticBag) {
	println("----TYPECHECKER----")
	symTable = table
	handler = bag
	for _, decl := range decls {
		checkDecl(decl)
	}
}
func areTypesEqual(pos error.Position, type1 *types.Type, type2 *types.Type) bool {
	if type1.TypeId == type2.TypeId {
		return true
	}
	handler.ReportError(pos, "Expected '%s' type but got '%s' type", type1.TypeName, type2.TypeName)
	return false
}
func checkDecl(decl ast.Decl) {
	switch node := decl.(type) {
	case *ast.DeclFunction:
		{
			enterScope(symTable.GetObj(node.Name).GetScope())
			functionName = node.Name
			for _, stmt := range node.Body.Block {
				checkStmt(stmt)
			}
			leaveScope()
		}
	}
}

func checkStmt(stmt ast.Stmt) {
	pos := stmt.GetPos()
	switch node := stmt.(type) {
	case *ast.StmtLet:
		{
			variableType := currScope.GetObj(node.Name).Type
			if node.Init != nil {
				result := checkExpr(node.Init, variableType)
				areTypesEqual(pos, variableType, result)
			}
		}
	case *ast.StmtIf:
		{
			cond := checkExpr(node.Cond, nil)
			if cond.Kind != types.TYPE_BOOL {
				handler.ReportError(pos, "if expression must be boolean")
			}
			checkStmt(node.Then)
		}
	case *ast.StmtBlock:
		{
			enterScope(node.Scope)
			for _, stmt := range node.Block {
				checkStmt(stmt)
			}
			leaveScope()
		}
	case *ast.StmtReturn:
		{
			returnType := symTable.GetObj(functionName).Type
			if returnType.Kind == types.TYPE_VOID && node.Result != nil {
				handler.ReportError(pos, "'%s' function shouldn't return anything", functionName)
				return
			}
			result := checkExpr(node.Result, returnType)
			areTypesEqual(pos, returnType, result)
		}
	case *ast.StmtExpr:
		{
			checkExpr(node.Expr, nil)
		}

	}
}
func checkExpr(expr ast.Expr, expectedType *types.Type) *types.Type {
	var typeResult *types.Type = nil
	switch node := expr.(type) {
	case *ast.ExprInt:
		{
			typeResult = symTable.GetObj("i8").Type
		}
	case *ast.ExprBoolean:
		{
			typeResult = symTable.GetObj("bool").Type
		}
	case *ast.ExprIdent:
		{
			return currScope.GetObj(node.Name).Type
		}
	case *ast.ExprAssign:
		{
			left := checkExpr(node.Left, expectedType)
			right := checkExpr(node.Right, expectedType)
			areTypesEqual(node.GetPos(), left, right)
		}
	case *ast.ExprGet:
		{
			obj := currScope.GetObj(node.Name)
			var typResult *types.Type
			enterScope(symTable.GetObj(obj.Type.TypeName).GetScope())
			typResult = checkExpr(node.Right, nil)
			leaveScope()
			return typResult
		}
	case *ast.ExprBinary:
		{
			left := checkExpr(node.Left, nil)
			right := checkExpr(node.Right, nil)
			// TODO: check if they are number or ptr
			switch node.Op {
			case token.TK_LESSEQUAL:
				fallthrough
			case token.TK_LESSTHAN:
				fallthrough
			case token.TK_GREATEREQUAL:
				fallthrough
			case token.TK_GREATERTHAN:
				fallthrough
			case token.TK_EQUAL:
				{
					if areTypesEqual(node.GetPos(), left, right) {
						return symTable.GetObj("bool").Type
					}
				}
			}
			return left
		}
	default:
		{
			panic("Unreachable")
		}
	}
	if typeResult != nil {
		if expectedType != nil && expectedType.Kind == typeResult.Kind {
			return expectedType
		}
		return typeResult
	}
	return symTable.GetObj("void").Type
}*/
