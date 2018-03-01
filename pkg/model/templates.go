package model

import (
	"fmt"

	"encoding/base64"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
)

func ValidateCreateTemplate(tmpl mttypes.Template) []error {
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

func ValidateUpdateTemplate(tmpl mttypes.Template) []error {
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
