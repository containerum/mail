package headers

import (
	"fmt"
)

//Extrernal headers
const (
	UserClientHeader    = "User-Client"
	AuthorizationHeader = "User-Token"
)

//Internal headers
const (
	RequestIDXHeader      = "X-Request-ID"
	RequestNameXHeader    = "X-Request-Name"
	UserIDXHeader         = "X-User-ID"
	UserClientXHeader     = "X-User-Client"
	UserAgentXHeader      = "X-User-Agent"
	UserIPXHeader         = "X-Client-IP"
	TokenIDXHeader        = "X-Token-ID"
	UserRoleXHeader       = "X-User-Role"
	UserNamespacesXHeader = "X-User-Namespace"
	UserVolumesXHeader    = "X-User-Volume"
	UserHideDataXHeader   = "X-User-Hide-Data"
)

//ErrHeaderRequired return if header not required
func ErrHeaderRequired(name string) error {
	return fmt.Errorf("Header %s required", name)
}

//ErrInvalidFormat return if header has invalid format
func ErrInvalidFormat(name string) error {
	return fmt.Errorf("Header %s has invalid format", name)
}
