package resource

import (
	"reflect"

	"regexp"

	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

type CreateResourceRequest struct {
	TariffID string `json:"tariff-id" binding:"uuid4"`
	Label    string `json:"label"`
}

type RenameResourceRequest struct {
	New string `json:"label"`
}

type SetResourceLockRequest struct {
	Lock bool `json:"lock" binding:"dns"`
}

type SetResourceAccessRequest struct {
	UserID string `json:"user_id" binding:"uuid4"`
	Access string `json:"access"`
}

// custom tag registration

func RegisterCustomTags(validate *validator.Validate) error {
	return validate.RegisterValidation("dns", dnsValidationFunc)
}

func RegisterCustomTagsGin(validate binding.StructValidator) error {
	return validate.RegisterValidation("dns", dnsValidationFunc)
}

var dnsLabel = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

func dnsValidationFunc(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	return dnsLabel.MatchString(field.String())
}
