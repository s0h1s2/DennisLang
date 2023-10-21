package error

import "fmt"

type Error struct {
	Pos Position
	Msg string
}

func (e Error) Error() string {
	return fmt.Sprintf("ERROR:(%d,%d):%s at line %d.", e.Pos.Start, e.Pos.End, e.Msg, e.Pos.Line)
}
