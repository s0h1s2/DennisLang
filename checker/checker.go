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
	symTable = table
	handler = bag
	for _, decl := range decls {
		checker(decl, nil)
	}
}
func checker(node ast.Node, expectedType *types.Type) *types.Type {
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
			resultType := exprChecker(n.Result, expectedType)
			if expectedType.Kind != resultType.Kind {
				handler.ReportError(error.Error{Msg: "function return type is wrong"})
			}
			return resultType
		}
	case *ast.StmtLet:
		{
			obj := symTable.GetObj(n.Name)
			if n.Init != nil {
				typ := exprChecker(n.Init, obj.Type)
				if obj.Type.TypeId != typ.TypeId {
					handler.ReportError(error.Error{Msg: "Types aren't equal"})
				}
			}
		}
	}

	return nil
}
func exprChecker(expr ast.Expr, expectedType *types.Type) *types.Type {
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
	case *ast.ExprBinary:
		{
			left := exprChecker(n.Left, expectedType)
			right := exprChecker(n.Right, expectedType)
			if right.Kind != types.TYPE_INT {
				handler.ReportError(error.Error{Msg: "Arithmetic operations for integers only"})
				return left
			}
			if left.Kind != types.TYPE_INT || left.TypeId != right.TypeId {
				handler.ReportError(error.Error{Msg: "Types aren't equal"})
			}
			return right
		}
	}
	if typeResult != nil {
		if typeResult.Kind == expectedType.Kind {
			return expectedType
		}
	}
	return symTable.GetObj("void").Type
}
