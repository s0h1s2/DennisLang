package resolver

import "github.com/s0h1s2/ast"

type Table struct {
	symbols map[string]*Object
}

func (t *Table) declare(name string) {
	if _, ok := t.symbols[name]; ok {
		// TODO
	}
}
func Resolve(ast []ast.Decl) *Table {
	table := Table{
		symbols: make(map[string]*Object),
	}
	for _, decl := range ast {
		resolver(decl)
	}
	return &table
}
func resolver(node ast.Node) {
	switch n := node.(type) {
	case *ast.DeclFunction:
		{

		}
	case *ast.StmtLet:
		{
		}

	}
}
