package error

import "fmt"

type Error struct {
	Pos Position
	Msg string
}

func (e Error) Error() string {
	return fmt.Sprintf("ERROR:%s at line %d", e.Msg, e.Pos.Line)
}
