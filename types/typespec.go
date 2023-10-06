package types

type TypeSpec interface {
	typeSpec()
}
type TypeName struct {
	Name string
}
type TypePtr struct {
	Base TypeSpec
}

func (ts *TypeName) typeSpec() {}
func (ts *TypePtr) typeSpec()  {}
