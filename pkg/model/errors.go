package model

import (
	"errors"

	"github.com/google/uuid"
)

const (
	isRequired      = "Field %v is required"
	IsRequiredQuery = "Field in query %v is required"
	notBase64       = "Field %v should be encoded in base64"
	moreZero        = "Field %v should be >0"
)

var (
	ErrInvalidID         = errors.New("ID should be UUID")
	ErrStorageOpenFailed = errors.New("Failed to open storage")
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	if err != nil {
		return false
	}
	return true
}
