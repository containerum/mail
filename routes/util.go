package routes

import (
	"fmt"

	"gopkg.in/go-playground/validator.v8"
)

//BindError is a type for bind errors
type BindError struct {
	Error string `json:"error"`
}

//ParseBindErorrs parses errors from message content binding
func ParseBindErorrs(in error) []BindError {
	var out []BindError

	t, isValidationError := in.(validator.ValidationErrors)

	if isValidationError {
		for _, v := range t {
			switch v.Tag {
			case "required":
				out = append(out, BindError{fmt.Sprintf("Field %v should be provided", v.Name)})
			default:
				out = append(out, BindError{fmt.Sprintf("%v should be %v", v.Name, v.Tag)})
			}
		}
		return out
	}
	return []BindError{{Error: in.Error()}}
}
