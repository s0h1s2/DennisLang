package c

import (
	"fmt"
	"os"
	"strings"

	"github.com/s0h1s2/resolver"
	"github.com/s0h1s2/scope"
)

type cCodegen struct {
	outputFile *os.File
	decls      []resolver.DeclNode
	table      *resolver.Table
}

func New(file *os.File, decls []resolver.DeclNode, table *resolver.Table) *cCodegen {
	return &cCodegen{
		outputFile: file,
		decls:      decls,
		table:      table,
	}
}
func (c *cCodegen) emitLine(format string, args ...any) {
	result := fmt.Sprintf(format, args...)
	c.outputFile.WriteString(result + "\n")
}
func (c *cCodegen) emit(format string, args ...any) {
	result := fmt.Sprintf(format, args...)
	c.outputFile.WriteString(result)
}
func (c *cCodegen) generateDecl(decl resolver.DeclNode) {
	switch node := decl.(type) {
	case *resolver.DeclFunction:
		{
			params := c.table.Symbols.GetObj(node.Name).Scope.QueryObjByKind(scope.PARAM)
			var sb strings.Builder
			paramsLength := len(params)
			for i, param := range params {
				sb.WriteString(fmt.Sprintf("%s %s", param.Type.TypeName, param.Name))
				if i < paramsLength-1 {
					sb.WriteByte(',')
				}
			}
			c.emit("%s %s(%s)", node.ReturnType.TypeName, node.Name, sb.String())
			c.generateStmt(node.Body)
			c.emitLine("")
		}
	case *resolver.DeclStruct:
		{
			c.emit("typedef struct %s{", node.Name)
			for _, field := range node.Fields {
				c.emit("%s %s;", field.Type.TypeName, field.Name)
			}
			c.emit("} %s;", node.Name)
			c.emit("\n")
		}
	}
}
func (c *cCodegen) generateStmt(stmt resolver.StmtNode) {
	switch node := stmt.(type) {
	case *resolver.StmtLet:
		{
			c.emit("%s %s", node.Type.TypeName, node.Name)
			if node.Init != nil {
				c.emit("=")
				c.generateExpr(node.Init)
			}
			c.emit(";")
		}
	case *resolver.StmtExpr:
		{
			c.generateExpr(node.Expr)
			c.emit(";")
		}
	case *resolver.StmtBlock:
		{
			c.emit("{")
			for _, stmt := range node.Body {
				c.generateStmt(stmt)
			}
			c.emitLine("}")
		}
	}
}
func (c *cCodegen) generateExpr(expr resolver.ExprNode) {
	switch node := expr.(type) {
	case *resolver.ExprInt:
		{
			c.emit(node.Value)
		}
	case *resolver.ExprAssign:
		{
			c.generateExpr(node.Left)
			c.emit("=")
			c.generateExpr(node.Right)
		}
	case *resolver.ExprCompound:
		{
			c.emit("(%s)", node.Type.TypeName)
			c.emit("{")
			for i, field := range node.Fields {
				c.emit(".%s=", field.Name)
				c.generateExpr(field.Expr)
				if i < len(node.Fields)-1 {
					c.emit(",")
				}
			}
			c.emit("}")
		}
	case *resolver.ExprField:
		{
			c.generateExpr(node.Expr)
			c.emit("." + node.Name)
		}
	case *resolver.ExprBool:
		{
			if node.Value {
				c.emit("true")
				return
			}
			c.emit("false")
		}
	case *resolver.ExprIdentifier:
		{
			c.emit(node.Name)
		}
	case *resolver.ExprString:
		{
			c.emit("\"%s\"", node.Value)
		}
	case *resolver.ExprCall:
		{
			c.emit(node.Name + "(")
			for i, arg := range node.Args {
				c.generateExpr(arg.Expr)
				if i < len(node.Args)-1 {
					c.emit(",")
				}
			}
			c.emit(")")
		}
	}
}
func (c *cCodegen) GenerateCode() {
	println("----CCODEGEN---")
	c.emitLine("// Declerations")
	c.emitLine("#include <stdio.h>")
	c.emitLine("#include <stdlib.h>")
	c.emitLine("#include <stdbool.h>")
	c.emitLine("typedef char i8;")
	c.emitLine("typedef short i16;")
	c.emitLine("typedef int i32;")
	c.emitLine("typedef long i64;")
	c.emitLine("typedef char* string;")
	for _, decl := range c.decls {
		c.generateDecl(decl)
	}
}
