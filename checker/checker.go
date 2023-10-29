package checker

import (
	"fmt"

	"github.com/s0h1s2/error"
	"github.com/s0h1s2/resolver"
	"github.com/s0h1s2/types"
)

type checker struct {
	handler    *error.DiagnosticBag
	symTable   *resolver.Table
	currentFun *resolver.DeclFunction
}

func Check(decls []resolver.DeclNode, table *resolver.Table, handler *error.DiagnosticBag) {
	println("----CHECKER----")
	c := &checker{handler: handler, symTable: table}
	for _, decl := range decls {
		c.checkDecl(decl)
	}
}
func (c *checker) areTypesEqual(type1 *types.Type, type2 *types.Type) bool {
	if type1 == nil || type2 == nil {
		return false
	}
	if type1.TypeId == type2.TypeId {
		return true
	}
	return false
}
func (c *checker) checkDecl(decl resolver.DeclNode) {
	switch node := decl.(type) {
	case *resolver.DeclFunction:
		{
			c.checkFunction(node)
		}
	case *resolver.DeclStruct:
		{
		}
	}
}
func (c *checker) checkFunction(fun *resolver.DeclFunction) {
	c.currentFun = fun
	c.checkStmt(fun.Body)
}
func (c *checker) checkStmt(stmt resolver.StmtNode) *types.Type {
	switch node := stmt.(type) {
	case *resolver.StmtBlock:
		{
			for _, stmt := range node.Body {
				c.checkStmt(stmt)
			}
		}
	case *resolver.StmtReturn:
		{
			if node.Result != nil {
				resultType := c.checkExpr(node.Result, c.currentFun.ReturnType)
				if !c.areTypesEqual(resultType, c.currentFun.ReturnType) {
					c.handler.ReportError(node.GetPos(), "Expected '%s' but got '%s' in function return", c.currentFun.ReturnType.TypeName, resultType.TypeName)
				}
			}
		}
	case *resolver.StmtExpr:
		{
			c.checkExpr(node.Expr, nil)
		}
	case *resolver.StmtLet:
		{
			if node.Init != nil {
				exprType := c.checkExpr(node.Init, node.Type)
				if !c.areTypesEqual(node.Type, exprType) {
					c.handler.ReportError(node.GetPos(), "Expected '%s' type but got '%s' type", node.Type.TypeName, exprType.TypeName)
				}
			}
		}
	}
	return nil
}

func (c *checker) checkExpr(expr resolver.ExprNode, expectedType *types.Type) *types.Type {
	var typeResult *types.Type = nil
	switch node := expr.(type) {
	case *resolver.ExprAssign:
		{
			left := c.checkExpr(node.Left, expectedType)
			right := c.checkExpr(node.Right, expectedType)
			if !c.areTypesEqual(left, right) {
				c.handler.ReportError(node.GetPos(), "Expected '%s' but got '%s'", left.TypeName, right.TypeName)
			}
			return left
		}
	case *resolver.ExprGet:
		{
			c.checkExpr(node.Right, node.Type)

		}
	case *resolver.ExprInt:
		{
			typeResult = c.symTable.Symbols.GetObj("i8").Type
		}
	case *resolver.ExprBool:
		{
			typeResult = c.symTable.Symbols.GetObj("bool").Type
		}
	case *resolver.ExprIdentifier:
		{
			println(node.Name)
			println(node.Type.TypeName)
			typeResult = node.Type
		}
	default:
		{
			panic(fmt.Sprintf("Unhandled %T or unreachable\n", node))
		}
	}
	if typeResult == nil {
		return c.symTable.Symbols.GetObj("void").Type
	}
	if expectedType == nil {
		return typeResult
	}
	if typeResult.Kind == expectedType.Kind {
		return expectedType
	}
	return typeResult
}
