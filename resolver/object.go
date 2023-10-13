package resolver

import (
	"github.com/s0h1s2/ast"
	"github.com/s0h1s2/types"
)

type ObjectKind int

const (
	FN ObjectKind = iota
	VAR
	TYPE
)

type Object struct {
	Kind ObjectKind
	Node ast.Node
	Type types.Type
}

func newObj(kind ObjectKind) *Object {
	return &Object{
		Kind: kind,
	}
}
