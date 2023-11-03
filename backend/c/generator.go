package c

import (
	"fmt"
	"os"

	"github.com/s0h1s2/resolver"
)

type cCodegen struct {
	outputFile *os.File
	decls      []resolver.DeclNode
}

func New(file *os.File, decls []resolver.DeclNode) *cCodegen {
	return &cCodegen{
		outputFile: file,
		decls:      decls,
	}
}
func (c *cCodegen) emitLine(format string, args ...any) {
	result := fmt.Sprintf(format, args...)
	c.outputFile.WriteString(result + "\n")
}
func (c *cCodegen) emit(line string) {
	c.outputFile.WriteString(line)
}
func (c *cCodegen) generateDecl(decl resolver.DeclNode) {
	switch node := decl.(type) {
	case *resolver.DeclFunction:
		{
			c.emitLine("%s %s(int argc,char **argv) {", node.ReturnType.TypeName, node.Name)
			c.emitLine("printf(\"Hello,from dennis compiler!\\n\");")
			c.emitLine("}")
		}
	case *resolver.DeclStruct:
		{

		}
	}
}
func (c *cCodegen) GenerateCode() {
	c.emitLine("#include <stdio.h>")
	c.emitLine("#include <stdlib.h>")
	c.emitLine("#define i8  char")
	c.emitLine("#define i16 short")
	c.emitLine("#define i32 int")
	c.emitLine("#define i64 long")
	for _, decl := range c.decls {
		c.generateDecl(decl)
	}
}
