package validation

import (
	"errors"

	"github.com/google/uuid"
)

const (
	isRequired = "Field %v is required"
	notBase64  = "Field %v should be encoded in base64"
	moreZero   = "Field %v should be >0"
)

//nolint
var (
	errInvalidID = errors.New("ID should be UUID")
)

// IsValidUUID checks if UUID is valid
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
