package routes

import (
	"encoding/json"
	"errors"
	"net"
	"regexp"

	"git.containerum.net/ch/grpc-proto-files/auth"
)

var uuidRegexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

var standardHeaderValidators = validators{
	"X-User-ID":         uuidValidator,
	"X-User-IP":         ipValidator,
	"X-User-Role":       roleValidator,
	"X-User-Part-Token": uuidValidator,
}

func uuidValidator(value string) error {
	if uuidRegexp.MatchString(value) {
		return nil
	} else {
		return errors.New("invalid UUID format")
	}
}

func ipValidator(value string) error {
	if ip := net.ParseIP(value); ip != nil {
		return nil
	} else {
		return errors.New("invalid IP format")
	}
}

func roleValidator(value string) error {
	if _, ok := auth.Role_value[value]; ok {
		return nil
	} else {
		return errors.New("invalid role")
	}
}

func resourcesAccessBodyValidator(body []byte) error {
	var bodyObj struct {
		Access *auth.ResourcesAccess `json:"access"`
	}

	if err := json.Unmarshal(body, &bodyObj); err == nil {
		return nil
	} else {
		return errors.New("invalid request body format")
	}
}
