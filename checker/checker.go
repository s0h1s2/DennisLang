package checker

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
			result := checkExpr(node.Init, variableType)
			areTypesEqual(pos, variableType, result)
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
}

//
// func checker(node ast.Node, expectedType *types.Type) *types.Type {
// 	pos := node.GetPos()
// 	switch n := node.(type) {
// 	case *ast.DeclFunction:
// 		{
// 			for _, stmt := range n.Body {
// 				checker(stmt, symTable.GetObj(n.Name).Type)
// 			}
// 			return nil
// 		}
// 	case *ast.StmtReturn:
// 		{
// 			if n.Result != nil && expectedType.Kind == types.TYPE_VOID {
// 				handler.ReportError(pos, "function shouldn't return becuase type is void")
// 				return nil
// 			}
// 			if n.Result == nil {
// 				handler.ReportError(pos, "Expected '%s' but got expression is empty", expectedType.TypeName)
// 				return nil
// 			}
// 			resultType := exprChecker(n.Result, expectedType)
// 			if expectedType.Kind != resultType.Kind {
// 				handler.ReportError(pos, "Expected '%s' type but got '%s' in return statement", expectedType.TypeName, resultType.TypeName)
// 			}
// 			return resultType
// 		}
// 	case *ast.StmtLet:
// 		{
// 			obj := symTable.GetObj(n.Name)
// 			if obj.Type.Kind == types.TYPE_VOID {
// 				handler.ReportError(pos, "Binding variable '%s' to 'void' type is not permitted", n.Name)
// 				return obj.Type
// 			}
// 			if n.Init != nil {
// 				typeResult := exprChecker(n.Init, obj.Type)
// 				if obj.Type.TypeId != typeResult.TypeId {
// 					handler.ReportError(pos, "Expected '%s' type but got '%s' type", obj.Type.TypeName, typeResult.TypeName)
// 				}
// 			}
// 		}
// 	}
//
// 	return nil
// }
// func exprChecker(expr ast.Expr, expectedType *types.Type) *types.Type {
// 	pos := expr.GetPos()
// 	var typeResult *types.Type = nil
// 	switch n := expr.(type) {
// 	case *ast.ExprInt:
// 		{
// 			typeResult = symTable.GetObj("i8").Type
// 		}
// 	case *ast.ExprBoolean:
// 		{
// 			typeResult = symTable.GetObj("bool").Type
// 		}
// 	case *ast.ExprIdent:
// 		{
// 			return symTable.GetObj(n.Name).Type
// 		}
// 	case *ast.ExprAssign:
// 		{
// 			left := exprChecker(n.Left, expectedType)
// 			right := exprChecker(n.Left, expectedType)
// 			if left.TypeId != right.TypeId {
// 				handler.ReportError(pos, "Expected '%s' but got '%s' when assign variable", left.TypeName, right.TypeName)
// 			}
// 		}
// 	case *ast.ExprAddrOf:
// 		{
// 			if expectedType.Kind == types.TYPE_PTR {
// 				typeResult := exprChecker(n.Right, expectedType.Base)
// 				if typeResult.TypeId != expectedType.Base.TypeId {
// 					handler.ReportError(pos, "Expected '%s' to point to '%s' but it is pointing to '%s' ", expectedType.TypeName, expectedType.Base.TypeName, typeResult.TypeName)
// 				}
// 				return expectedType
// 			}
// 		}
// 	case *ast.ExprBinary:
// 		{
// 			left := exprChecker(n.Left, expectedType)
// 			right := exprChecker(n.Right, expectedType)
// 			if left.Kind != types.TYPE_INT {
// 				handler.ReportError(pos, "Only integers can be '+' or '*'")
// 				return left
// 			}
// 			if left.TypeId != right.TypeId {
// 				handler.ReportError(pos, "Expected '%s' but go '%s' ", left.TypeName, right.TypeName)
// 			}
// 			return right
// 		}
// 	}
// 	if typeResult != nil {
// 		if typeResult.Kind == expectedType.Kind {
// 			return expectedType
// 		}
// 		return typeResult
// 	}
// 	return symTable.GetObj("void").Type
// }
