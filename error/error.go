package error

import "fmt"

type Error struct {
	Pos Position
	Msg string
}

func (e Error) Error() string {
	if e.Pos.isValid() {
		return fmt.Sprintf("ERROR:(%d,%d):%s at line %d.", e.Pos.Start, e.Pos.End, e.Msg, e.Pos.Line)
	}
	return fmt.Sprintf("ERROR:%s.", e.Msg)
}
