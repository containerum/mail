package httputil

import (
	"context"
	"fmt"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
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

var headersKey = new(struct{})

// SaveHeaders is a gin middleware which saves headers to request context
func SaveHeaders(ctx *gin.Context) {
	rctx := context.WithValue(ctx.Request.Context(), headersKey, ctx.Request.Header)
	ctx.Request = ctx.Request.WithContext(rctx)
}

// RequestHeadersMap extracts saved headers from context as map[string]string (useful for resty).
// saveHeaders middleware required for operation.
func RequestHeadersMap(ctx context.Context) map[string]string {
	ret := make(map[string]string)
	for k, v := range ctx.Value(headersKey).(http.Header) {
		if len(v) > 0 {
			ret[textproto.CanonicalMIMEHeaderKey(k)] = v[0] // this is how MIMEHeader.Get() works actually
		}
	}
	return ret
}

// RequestXHeadersMap works like RequestHeadersMap but returns only "X-" headers
func RequestXHeadersMap(ctx context.Context) map[string]string {
	ret := make(map[string]string)
	for k, v := range ctx.Value(headersKey).(http.Header) {
		k = textproto.CanonicalMIMEHeaderKey(k)
		if len(v) > 0 && strings.HasPrefix(k, "X-") {
			ret[k] = v[0]
		}
	}
	return ret
}

// RequestHeaders extracts saved headers from context.
// saveHeaders middleware required for operation.
func RequestHeaders(ctx context.Context) http.Header {
	return ctx.Value(headersKey).(http.Header)
}

var hdrToKey = map[string]interface{}{
	textproto.CanonicalMIMEHeaderKey(UserIDXHeader):     UserIDContextKey,
	textproto.CanonicalMIMEHeaderKey(UserAgentXHeader):  UserAgentContextKey,
	textproto.CanonicalMIMEHeaderKey(UserClientXHeader): FingerPrintContextKey,
	textproto.CanonicalMIMEHeaderKey(RequestIDXHeader):  RequestIDContextKey,
	textproto.CanonicalMIMEHeaderKey(TokenIDXHeader):    TokenIDContextKey,
	textproto.CanonicalMIMEHeaderKey(UserIPXHeader):     ClientIPContextKey,
	textproto.CanonicalMIMEHeaderKey(UserRoleXHeader):   UserRoleContextKey,
}

// RequireHeaders is a gin middleware to ensure that headers is set
func RequireHeaders(errToReturn cherry.ErrConstruct, headers ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var notFoundHeaders []string
		for _, v := range headers {
			if ctx.GetHeader(textproto.CanonicalMIMEHeaderKey(v)) == "" {
				notFoundHeaders = append(notFoundHeaders, v)
			}
		}
		if len(notFoundHeaders) > 0 {
			err := errToReturn()
			for _, notFoundHeader := range notFoundHeaders {
				err.AddDetailF("required header %s was not provided", notFoundHeader)
			}
			gonic.Gonic(err, ctx)
		}
	}
}

// PrepareContext is a gin middleware which adds values from header to context
func PrepareContext(ctx *gin.Context) {
	for hn, ck := range hdrToKey {
		if hv := ctx.GetHeader(hn); hv != "" {
			rctx := context.WithValue(ctx.Request.Context(), ck, hv)
			ctx.Request = ctx.Request.WithContext(rctx)
		}
	}

	acceptLanguages := ctx.GetHeader("Accept-Language")
	acceptLanguagesToContext := make([]string, 0)
	for _, language := range strings.Split(acceptLanguages, ",") {
		language = strings.Split(strings.TrimSpace(language), ";")[0] // drop quality values
		acceptLanguagesToContext = append(acceptLanguagesToContext, language)
	}
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), AcceptLanguageContextKey, acceptLanguagesToContext))
}

// RequireAdminRole is a gin middleware which requires admin role
func RequireAdminRole(errToReturn cherry.ErrConstruct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader(textproto.CanonicalMIMEHeaderKey(UserRoleXHeader)) != "admin" {
			err := errToReturn().AddDetails("only admin can do this")
			gonic.Gonic(err, ctx)
		}
	}
}

// SubstituteUserMiddleware replaces user id in context with user id from query if it set and user is admin
func SubstituteUserMiddleware(validate *validator.Validate, translator *ut.UniversalTranslator, validationErr cherry.ErrConstruct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetHeader(textproto.CanonicalMIMEHeaderKey(UserRoleXHeader))
		if userID, set := ctx.GetQuery("user-id"); set && role == "admin" {
			if vErr := validate.VarCtx(ctx.Request.Context(), userID, "uuid"); vErr != nil {
				t, _ := translator.FindTranslator(GetAcceptedLanguages(ctx.Request.Context())...)
				err := validationErr().AddDetailF("Parameter \"user-id\": %s", vErr.(validator.ValidationErrors).Translate(t))
				gonic.Gonic(err, ctx)
				return
			}
			rctx := context.WithValue(ctx.Request.Context(), UserIDContextKey, userID)
			ctx.Request = ctx.Request.WithContext(rctx)
		}
	}
}
