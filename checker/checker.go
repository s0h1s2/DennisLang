package checker

import (
	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/resolver"
	"github.com/s0h1s2/types"
)

var symTable *resolver.Table
var declerations []ast.Decl
var handler *error.DiagnosticBag

func TypeChecker(table *resolver.Table, decls []ast.Decl, bag *error.DiagnosticBag) {
	println("----TYPECHECKER----")
	symTable = table
	handler = bag
	for _, decl := range decls {
		checker(decl, nil)
	}
}
func checker(node ast.Node, expectedType *types.Type) *types.Type {
	pos := node.GetPos()
	switch n := node.(type) {
	case *ast.DeclFunction:
		{
			for _, stmt := range n.Body {
				checker(stmt, symTable.GetObj(n.Name).Type)
			}
			return nil
		}
	case *ast.StmtReturn:
		{
			if n.Result != nil && expectedType.Kind == types.TYPE_VOID {
				handler.ReportError(pos, "function shouldn't return becuase type is void")
				return nil
			}
			if n.Result == nil {
				handler.ReportError(pos, "Expected '%s' but got expression is empty", expectedType.TypeName)
				return nil
			}
			resultType := exprChecker(n.Result, expectedType)
			if expectedType.Kind != resultType.Kind {
				handler.ReportError(pos, "Expected '%s' type but got '%s' in return statement", expectedType.TypeName, resultType.TypeName)
			}
			return resultType
		}
	case *ast.StmtLet:
		{
			obj := symTable.GetObj(n.Name)
			if obj.Type.Kind == types.TYPE_VOID {
				handler.ReportError(pos, "Binding variable '%s' to 'void' type is not permitted", n.Name)
				return obj.Type
			}
			if n.Init != nil {
				typeResult := exprChecker(n.Init, obj.Type)
				if obj.Type.TypeId != typeResult.TypeId {
					handler.ReportError(pos, "Expected '%s' type but got '%s' type", expectedType.TypeName, typeResult.TypeName)
				}
			}
		}
	}

	return nil
}
func exprChecker(expr ast.Expr, expectedType *types.Type) *types.Type {
	pos := expr.GetPos()
	var typeResult *types.Type = nil
	switch n := expr.(type) {
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
			return symTable.GetObj(n.Name).Type
		}
	case *ast.ExprAssign:
		{
			left := exprChecker(n.Left, expectedType)
			right := exprChecker(n.Left, expectedType)
			if left.TypeId != right.TypeId {
				handler.ReportError(pos, "Expected '%s' but got '%s' when assign variable", left.TypeName, right.TypeName)
			}
		}
	case *ast.ExprBinary:
		{
			left := exprChecker(n.Left, expectedType)
			right := exprChecker(n.Right, expectedType)
			if left.Kind != types.TYPE_INT {
				handler.ReportError(pos, "Only integers can be '+' or '*'")
				return left
			}
			if left.TypeId != right.TypeId {
				handler.ReportError(pos, "Expected '%s' but go '%s' ", left.TypeName, right.TypeName)
			}
			return right
		}
	}
	if typeResult != nil {
		if typeResult.Kind == expectedType.Kind {
			return expectedType
		}
		return typeResult
	}
	return symTable.GetObj("void").Type
}
