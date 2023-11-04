package types

var typeId int = 0

type TypeKind int

const (
	TYPE_INT TypeKind = iota
	TYPE_VOID
	TYPE_NULL
	TYPE_BOOL
	TYPE_PTR
	TYPE_STRUCT
)

type Type struct {
	TypeName  string
	Kind      TypeKind
	Size      uint64
	Alignment uint64
	Base      *Type
	TypeId    int
}

func NewType(name string, kind TypeKind, size uint64, align uint64) *Type {
	typeId++
	return &Type{
		TypeName:  name,
		Kind:      kind,
		Size:      size,
		Alignment: align,
		TypeId:    typeId,
		Base:      nil,
	}
}
