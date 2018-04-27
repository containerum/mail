package httputil

import "context"

type contextKey int

// Keys to inject data to context
const (
	FingerPrintContextKey contextKey = iota
	ClientIPContextKey
	UserAgentContextKey
	RequestIDContextKey
	UserIDContextKey
	TokenIDContextKey
	UserRoleContextKey

	AcceptLanguageContextKey
)

// MustGetFingerprint attempts to extract client fingerprint using FingerPrintContextKey from context.
// It panics if value was not found.
func MustGetFingerprint(ctx context.Context) string {
	fp, ok := ctx.Value(FingerPrintContextKey).(string)
	if !ok {
		panic("fingerprint not found in context")
	}
	return fp
}

// MustGetClientIP attempts to extract client IP address using ClientIPContextKey from context.
// It panics if value was not found.
func MustGetClientIP(ctx context.Context) string {
	ip, ok := ctx.Value(ClientIPContextKey).(string)
	if !ok {
		panic("client ip not found in context")
	}
	return ip
}

// MustGetUserAgent attempts to extract client IP address using UserAgentContextKey from context.
// It panics if value was not found.
func MustGetUserAgent(ctx context.Context) string {
	ip, ok := ctx.Value(UserAgentContextKey).(string)
	if !ok {
		panic("user agent not found in context")
	}
	return ip
}

// MustGetSessionID attempts to extract session ID using RequestIDContextKey from context.
// It panics if value was not found in context.
func MustGetRequestID(ctx context.Context) string {
	sid, ok := ctx.Value(RequestIDContextKey).(string)
	if !ok {
		panic("session id not found in context")
	}
	return sid
}

// MustGetUserID attempts to extract user ID using RequestIDContextKey from context.
// It panics if value was not found in context.
func MustGetUserID(ctx context.Context) string {
	uid, ok := ctx.Value(UserIDContextKey).(string)
	if !ok {
		panic("user id not found in context")
	}
	return uid
}

// MustGetTokenID attempts to extract token ID using TokenIDContextKey from context.
// It panics if value was not found in context.
func MustGetTokenID(ctx context.Context) string {
	uid, ok := ctx.Value(TokenIDContextKey).(string)
	if !ok {
		panic("token id not found in context")
	}
	return uid
}

// MustGetUserRole attempts to extract user role using UserRoleContextKey from context
// It panics if value was not found in context.
func MustGetUserRole(ctx context.Context) string {
	role, ok := ctx.Value(UserRoleContextKey).(string)
	if !ok {
		panic("user role not found in context")
	}
	return role
}

// GetAcceptedLanguages extracts accepted languages from context
func GetAcceptedLanguages(ctx context.Context) []string {
	alangs, _ := ctx.Value(AcceptLanguageContextKey).([]string)
	return alangs
}
