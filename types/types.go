package types

type Type struct {
	Size      uint64
	Alignment uint64
	Base      *Type
}
