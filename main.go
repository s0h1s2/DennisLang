package main

import (
	"fmt"
	"io"
	"os"

	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/checker"
	"github.com/s0h1s2/error"
	"github.com/s0h1s2/lexer"
	"github.com/s0h1s2/parser"
	"github.com/s0h1s2/resolver"
)

func printAst(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.ExprAssign:
		{
			printAst(e.Left)
			print("=")
			printAst(e.Right)
		}
	case *ast.ExprInt:
		{
			println(e.Value)
		}

	case *ast.ExprBinary:
		{
			printAst(e.Left)
			printAst(e.Right)
		}

	case *ast.ExprIdent:
		{
			println(e.Name)
		}
	default:
		{
			panic("UnReachable")
		}
	}
}
func hello(n ast.Node) bool {
	println(n)
	return false
}
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s path-to-file\n", os.Args[0])
		return
	}
	filePath := os.Args[1]
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("Provided file '%s' doesn't exist.", filePath)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		println("Unable to open file")
		return
	}
	src, err := io.ReadAll(file)
	file.Close()
	bag := error.New()
	lex := lexer.New(bag)
	tokens := lex.GetTokens([]byte(src))
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	// Passes
	parser := parser.New(tokens, bag)
	tree := parser.Parse()
	if bag.GotErrors() {
		bag.PrintErrors()
		return
	}

	table, resolvedDecls := resolver.Resolve(tree, bag)
	if bag.GotErrors() {
		bag.PrintErrors()
		return
	}
	checker.Check(resolvedDecls, table, bag)
	if bag.GotErrors() {
		bag.PrintErrors()
		return
	}

}
