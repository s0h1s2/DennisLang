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
	// Node  ast.Node
	Type  *types.Type
	Scope *Scope
}

func NewObj(kind ObjectKind, typee *types.Type, scope *Scope) *Object {
	return &Object{
		Kind:  kind,
		Type:  typee,
		Scope: scope,
	}
}
func NewTypeObj(typee *types.Type) *Object {
	return &Object{
		Kind: TYPE,
		Type: typee,
	}
}
func (o *Object) GetScope() *Scope {
	return o.Scope
}
