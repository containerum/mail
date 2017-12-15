package utils

import "fmt"

type Error struct {
	Text string `json:"error"`
}

func (e *Error) Error() string {
	return e.Text
}

func NewError(text string) *Error {
	return &Error{
		Text: text,
	}
}

func NewErrorF(format string, data ...interface{}) *Error {
	return &Error{
		Text: fmt.Sprintf(format, data...),
	}
}