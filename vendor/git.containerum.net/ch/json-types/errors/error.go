package errors

import "fmt"

type Error struct {
	Text string `json:"error"`
	Code int    `json:"code,omitempty"`
}

func (e *Error) Error() string {
	if e.Code == 0 {
		return e.Text
	}
	return fmt.Sprintf("description: %s, code: %d", e.Text, e.Code)
}

func New(text string) *Error {
	return &Error{
		Text: text,
	}
}

func NewWithCode(text string, code int) *Error {
	return &Error{
		Text: text,
		Code: code,
	}
}

func Format(format string, data ...interface{}) *Error {
	return &Error{
		Text: fmt.Sprintf(format, data...),
	}
}

func FormatWithCode(code int, format string, data ...interface{}) *Error {
	return &Error{
		Text: fmt.Sprintf(format, data...),
		Code: code,
	}
}
