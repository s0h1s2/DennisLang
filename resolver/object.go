package resolver

import "github.com/s0h1s2/ast"

type ObjectKind int

const (
	FN ObjectKind = iota
	VAR
)

type Object struct {
	Kind ObjectKind
	Decl ast.Decl
	//typ Types.
}

func newObj(kind ObjectKind) *Object {
	return &Object{
		Kind: kind,
	}
}
