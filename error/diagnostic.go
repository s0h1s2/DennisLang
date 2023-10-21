package error

import (
	"fmt"
)

type DiagnosticBag struct {
	errors []Error
}

func New() *DiagnosticBag {
	return &DiagnosticBag{
		errors: make([]Error, 0, 4),
	}
}

func (bag *DiagnosticBag) ReportError(pos Position, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	bag.errors = append(bag.errors, Error{Msg: msg, Pos: pos})
}
func (bag DiagnosticBag) PrintErrors() {
	colorReset := "\033[0m"
	colorRed := "\033[35m"
	print(colorRed)
	for _, err := range bag.errors {
		println(err.Error())
	}
	print(colorReset)
}
func (bag DiagnosticBag) GotErrors() bool {
	return len(bag.errors) > 0
}
