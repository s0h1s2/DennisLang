package error

type Position struct {
	Start int
	End   int
	Line  int
}

func (p Position) isValid() bool {
	return p.Line > 0
}
