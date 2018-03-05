package mtErrors

import (
	bytes "bytes"
	template "text/template"

	cherry "git.containerum.net/ch/kube-client/pkg/cherry"
)

const ()

// ErrAdminRequired error
// User is not admin and has no permissions
func ErrAdminRequired(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Admin access required", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x7, Kind: 0x1}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrPermissionsError error
// Unable to verify if user has required permissions, so don't give access
func ErrPermissionsError(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to verify permissions", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x7, Kind: 0x2}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrRequiredHeadersNotProvided error
// Required headers is not found in context
func ErrRequiredHeadersNotProvided(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Required headers not provided", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x7, Kind: 0x3}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrRequestValidationFailed error
// Validation error when parsing request
func ErrRequestValidationFailed(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Request validation failed", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x7, Kind: 0x4}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetTemplatesList(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get templates list", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0x5}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get template", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0x6}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableSaveTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to save template", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0x7}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableUpdateTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to update template", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0x8}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableDeleteTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to delete template", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0x9}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrTemplateAlreadyExists(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Template with this name already exists", StatusHTTP: 409, ID: cherry.ErrID{SID: 0x7, Kind: 0xa}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrTemplateNotExist(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Template with this name doesn't exist", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x7, Kind: 0xb}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrTemplateVersionNotExist(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Template with this name and version doesn't exist", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x7, Kind: 0xc}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetMessagesList(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get messages list", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0xd}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetMessage(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get message", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0xe}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableSaveMessage(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to save message", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0xf}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrMessageNotExist(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Message with this name doesn't exist", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x7, Kind: 0x10}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrMailSendFailed(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to send email", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x7, Kind: 0x11}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}
func renderTemplate(templText string) string {
	buf := &bytes.Buffer{}
	templ, err := template.New("").Parse(templText)
	if err != nil {
		return err.Error()
	}
	err = templ.Execute(buf, map[string]string{})
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
