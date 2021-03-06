package validation

import (
	"fmt"

	"encoding/base64"

	"git.containerum.net/ch/mail-templater/pkg/models"
)

//ValidateCreateTemplate validates template creation request
//nolint: gocyclo
func ValidateCreateTemplate(tmpl models.Template) []error {
	errs := []error{}
	if tmpl.Name == "" {
		errs = append(errs, fmt.Errorf(isRequired, "Name"))
	}
	if tmpl.Version == "" {
		errs = append(errs, fmt.Errorf(isRequired, "Version"))
	}
	if tmpl.Subject == "" {
		errs = append(errs, fmt.Errorf(isRequired, "Subject"))
	}
	if tmpl.Data == "" {
		errs = append(errs, fmt.Errorf(isRequired, "Data"))
	} else {
		_, err := base64.StdEncoding.DecodeString(tmpl.Data)
		if err != nil {
			errs = append(errs, fmt.Errorf(notBase64, "Data"))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

//ValidateUpdateTemplate validates template update request
//nolint: gocyclo
func ValidateUpdateTemplate(tmpl models.Template) []error {
	errs := []error{}
	if tmpl.Subject == "" {
		errs = append(errs, fmt.Errorf(isRequired, "Subject"))
	}
	if tmpl.Data == "" {
		errs = append(errs, fmt.Errorf(isRequired, "Data"))
	} else {
		_, err := base64.StdEncoding.DecodeString(tmpl.Data)
		if err != nil {
			errs = append(errs, fmt.Errorf(notBase64, "Data"))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
