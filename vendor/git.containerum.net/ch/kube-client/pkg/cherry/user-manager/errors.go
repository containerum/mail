package umErrors

import (
	bytes "bytes"
	template "text/template"

	cherry "git.containerum.net/ch/kube-client/pkg/cherry"
)

const ()

// ErrAdminRequired error
// User is not admin and has no permissions
func ErrAdminRequired(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Admin access required", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x4, Kind: 0x1}, Details: []string(nil)}
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
	err := &cherry.Err{Message: "Unable to verify permissions", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x4, Kind: 0x2}, Details: []string(nil)}
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
	err := &cherry.Err{Message: "Required headers not provided", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x4, Kind: 0x3}, Details: []string(nil)}
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
	err := &cherry.Err{Message: "Request validation failed", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x4, Kind: 0x4}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrInvalidRecaptcha error
// Invalid Recaptcha. Or you are a robot
func ErrInvalidRecaptcha(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Recaptcha validation failed", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x4, Kind: 0x5}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrInvalidLogin error
// Invalid login/password/token
func ErrInvalidLogin(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Invalid login credentials", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x4, Kind: 0x6}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrLoginFailed error
// Unable to login
func ErrLoginFailed(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Login failed", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x7}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrLogoutFailed error
// Unable to logout
func ErrLogoutFailed(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Logout failed", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x8}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrNotActivated error
// User should activate his account first
func ErrNotActivated(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Account is not activated", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x4, Kind: 0x9}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableActivate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to activate user account", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0xa}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableResetPassword(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to reset password", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0xb}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableChangePassword(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to change password", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0xc}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetUsersList(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get users list", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0xd}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetUserInfo(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get user", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0xe}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableCreateUser(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to create user", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0xf}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableUpdateUserInfo(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to update user", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x10}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableDeleteUser(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to delete user", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x11}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUserAlreadyExists(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "User with such credentials already exists", StatusHTTP: 409, ID: cherry.ErrID{SID: 0x4, Kind: 0x12}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUserNotExist(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "User with such credentials doesn't exist", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x4, Kind: 0x13}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetUserLinks(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get user links", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x14}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableBindAccount(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to bind account", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x15}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableUnbindAccount(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to unbind account", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x16}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableBlacklistUser(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to add user to blacklist", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x17}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableUnblacklistUser(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to remove user from blacklist", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x18}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUserAlreadyBlacklisted(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "User is already blacklisted", StatusHTTP: 409, ID: cherry.ErrID{SID: 0x4, Kind: 0x19}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUserNotBlacklisted(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "User is not blacklisted", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x4, Kind: 0x1a}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableBlacklistDomain(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "User is already blacklisted", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x1b}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableUnblacklistDomain(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "User is not blacklisted", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x1c}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrDomainAlreadyBlacklisted(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "User is already blacklisted", StatusHTTP: 409, ID: cherry.ErrID{SID: 0x4, Kind: 0x1d}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrDomainNotBlacklisted(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "User is not blacklisted", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x4, Kind: 0x1e}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableResendLink(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to resend link", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x1f}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInvalidLink(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Invalid link", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x4, Kind: 0x20}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetUserBlacklist(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get blacklisted users", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x21}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetDomainBlacklist(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get blacklisted domains", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x22}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInternalError(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Internal error", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x4, Kind: 0x23}, Details: []string(nil)}
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
