package validation

import (
	"fmt"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"github.com/goware/emailx"
)

//ValidateSimpleSendRequest validates simple send mail request
//nolint: gocyclo
func ValidateSimpleSendRequest(snd models.SimpleSendRequest) []error {
	errs := []error{}
	if snd.UserID == "" {
		errs = append(errs, fmt.Errorf(isRequired, "UserID"))
	} else if !IsValidUUID(snd.UserID) {
		errs = append(errs, errInvalidID)
	}
	if snd.Template == "" {
		errs = append(errs, fmt.Errorf(isRequired, "Template"))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

//ValidateSendRequest validates send mail request
//nolint: gocyclo
func ValidateSendRequest(snd models.SendRequest) []error {
	errs := []error{}
	if snd.Delay < 0 {
		errs = append(errs, fmt.Errorf(moreZero, "Delay"))
	}
	if snd.Message.Recipients == nil || len(snd.Message.Recipients) == 0 {
		errs = append(errs, fmt.Errorf(isRequired, "Recipients"))
	} else {
		for _, v := range snd.Message.Recipients {
			if v.Name == "" {
				errs = append(errs, fmt.Errorf(isRequired, "Name"))
			}
			if v.ID == "" {
				errs = append(errs, fmt.Errorf(isRequired, "ID"))
			} else if !IsValidUUID(v.ID) {
				errs = append(errs, errInvalidID)
			}
			if v.Email == "" {
				errs = append(errs, fmt.Errorf(isRequired, "Email"))
			} else if err := emailx.ValidateFast(v.Email); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
