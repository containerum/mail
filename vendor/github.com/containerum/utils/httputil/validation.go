package httputil

import (
	"strings"

	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func mapValues(m map[string]string) (ret []string) {
	for _, v := range m {
		ret = append(ret, v)
	}
	return
}

// ValidateQueryParamsMiddleware validates query parameters with provided tags. Key of "vmap" is parameter name, value is tag
func ValidateQueryParamsMiddleware(vmap map[string]string, validate *validator.Validate, translator *ut.UniversalTranslator, validationErr cherry.ErrConstruct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errMap := make(map[string]validator.ValidationErrorsTranslations)
		t, _ := translator.FindTranslator(GetAcceptedLanguages(ctx.Request.Context())...)
		for param, tag := range vmap {
			for _, value := range ctx.Request.URL.Query()[param] {
				vErr := validate.VarCtx(ctx.Request.Context(), value, tag)
				if vErr != nil {
					errMap[param] = vErr.(validator.ValidationErrors).Translate(t)
				}
			}
		}
		if len(errMap) > 0 {
			retErr := validationErr()
			for param, err := range errMap {
				retErr.AddDetailF("Query parameter \"%s\": %s", param, strings.Join(mapValues(err), ", "))
			}
			gonic.Gonic(retErr, ctx)
			return
		}
	}
}

// ValidateURLParamsMiddleware validates URL parameters with provided tags. Key of "vmap" is parameter name, value is tag
func ValidateURLParamsMiddleware(vmap map[string]string, validate *validator.Validate, translator *ut.UniversalTranslator, validationErr cherry.ErrConstruct) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errMap := make(map[string]validator.ValidationErrorsTranslations)
		t, _ := translator.FindTranslator(GetAcceptedLanguages(ctx.Request.Context())...)
		for param, tag := range vmap {
			vErr := validate.VarCtx(ctx.Request.Context(), ctx.Param(param), tag)
			if vErr != nil {
				errMap[param] = vErr.(validator.ValidationErrors).Translate(t)
			}
		}
		if len(errMap) > 0 {
			retErr := validationErr()
			for param, err := range errMap {
				retErr.AddDetailF("URL parameter \"%s\": %s", param, strings.Join(mapValues(err), ", "))
			}
			gonic.Gonic(retErr, ctx)
			return
		}
	}
}
