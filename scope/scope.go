package scope

type Scope struct {
	parent  *Scope
	symbols map[string]*Object
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: make(map[string]*Object),
	}
}
func (s *Scope) Lookup(name string) bool {
	scope := s
	for scope != nil {
		if _, ok := scope.symbols[name]; ok {
			return true
		}
		scope = scope.parent
	}
	return false
}
func (s *Scope) LookupOnce(name string) bool {
	if _, ok := s.symbols[name]; ok {
		return true
	}
	return false
}

func (s *Scope) Define(name string, obj *Object) *Object {
	s.symbols[name] = obj
	return obj
}
func (s *Scope) QueryByKind(kind ObjectKind) []string {
	objects := make([]string, 0, 4)
	for k, v := range s.symbols {
		if v.Kind == kind {
			objects = append(objects, k)
		}
	}
	return objects
}
func (s *Scope) QueryObjByKind(kind ObjectKind) []*Object {
	objects := make([]*Object, 0, 4)
	for _, v := range s.symbols {
		if v.Kind == kind {
			objects = append(objects, v)
		}
	}
	return objects
}

func (s *Scope) GetObj(name string) *Object {
	scope := s
	for scope != nil {
		if scope.LookupOnce(name) {
			return scope.symbols[name]
		}
		scope = scope.parent
	}
	return nil
}
