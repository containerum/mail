package middleware

import (
	"net/textproto"

	"git.containerum.net/ch/api-gateway/pkg/utils/headers"
	"git.containerum.net/ch/cherry/adaptors/gonic"
	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
	"github.com/gin-gonic/gin"
)

//RequireAdminRole checks if user is admin
func RequireAdminRole(ctx *gin.Context) {
	if ctx.GetHeader(textproto.CanonicalMIMEHeaderKey(headers.UserRoleXHeader)) != "admin" {
		gonic.Gonic(mtErrors.ErrAdminRequired(), ctx)
		return
	}
}
