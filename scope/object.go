package scope

import (
	// "github.com/s0h1s2/ast"
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
	Type *types.Type
}

func NewObj(kind ObjectKind, typee *types.Type) *Object {
	return &Object{
		Kind: kind,
		Type: typee,
	}
}
func NewTypeObj(typee *types.Type) *Object {
	return &Object{
		Kind: TYPE,
		Type: typee,
	}
}
