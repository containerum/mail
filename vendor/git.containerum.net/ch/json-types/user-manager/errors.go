package user

import "fmt"

type WebAPIError struct {
	Message interface{} `json:"message"` // can be also object
}

func (e *WebAPIError) Error() string {
	return fmt.Sprintf("%#v", e.Message)
}
