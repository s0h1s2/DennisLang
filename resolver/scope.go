package resolver

type Scope struct {
	parent  *Scope
	symbols map[string]*Object
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: make(map[string]*Object, 4),
	}
}
func (s *Scope) Lookup(name string) bool {
	scope := s
	for scope != nil {
		if _, ok := s.symbols[name]; ok {
			return true
		}
		scope = scope.parent
	}
	return false
}
func (s *Scope) Define(name string, obj *Object) bool {
	if ok := s.Lookup(name); !ok {
		s.symbols[name] = obj
		return true
	}
	return false
}
func (s *Scope) GetObj(name string) *Object {
	if s.Lookup(name) {
		return s.symbols[name]
	}
	return nil
}
