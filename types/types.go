package types

var typeId int = 0

type Type struct {
	Size      uint64
	Alignment uint64
	Base      *Type
	TypeId    int
}

func NewType(size uint64, align uint64) *Type {
	typeId++
	return &Type{
		Size:      size,
		Alignment: align,
		TypeId:    typeId,
		Base:      nil,
	}
}
