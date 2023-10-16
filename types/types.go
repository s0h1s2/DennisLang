package types

var typeId int = 0

type TypeKind int

const (
	TYPE_INT TypeKind = iota
	TYPE_VOID
	TYPE_BOOL
	TYPE_PTR
)

type Type struct {
	Kind      TypeKind
	Size      uint64
	Alignment uint64
	Base      *Type
	TypeId    int
}

func NewType(kind TypeKind, size uint64, align uint64) *Type {
	typeId++
	return &Type{
		Kind:      kind,
		Size:      size,
		Alignment: align,
		TypeId:    typeId,
		Base:      nil,
	}
}
