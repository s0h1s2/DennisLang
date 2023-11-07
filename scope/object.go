package scope

import (
	// "github.com/s0h1s2/ast"
	"github.com/s0h1s2/types"
)

type ObjectKind int

const (
	FN ObjectKind = iota
	VAR
	PARAM
	FIELD
	TYPE
)

type Object struct {
	Name  string
	Kind  ObjectKind
	Type  *types.Type
	Scope *Scope
}

func NewObj(name string, kind ObjectKind, typee *types.Type) *Object {
	return &Object{
		Name: name,
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
func NewFuncObj(params map[string]*types.Type, returnType *types.Type) *Object {
	fnScope := NewScope(nil)
	for name, obj := range params {
		fnScope.Define(name, NewObj(name, PARAM, obj))
	}
	return &Object{
		Kind:  FN,
		Type:  returnType,
		Scope: fnScope,
	}
}
