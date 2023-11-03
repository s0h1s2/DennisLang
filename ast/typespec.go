package ast

import "github.com/s0h1s2/error"

type TypeSpec interface {
	typeSpec()
	GetPos() error.Position
}
type TypeName struct {
	Name string
	Pos  error.Position
}
type TypePtr struct {
	Pos  error.Position
	Base TypeSpec
}

func (ts *TypeName) typeSpec() {}
func (ts *TypeName) GetPos() error.Position {
	return ts.Pos
}

func (ts *TypePtr) typeSpec() {}
func (ts *TypePtr) GetPos() error.Position {
	return ts.Pos
}
