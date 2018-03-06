package middleware

import (
	umtypes "git.containerum.net/ch/json-types/user-manager"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/mail-templater"
	"github.com/gin-gonic/gin"
)

//RequireAdminRole checks if user is admin
func RequireAdminRole(ctx *gin.Context) {
	if ctx.GetHeader(umtypes.UserRoleHeader) != "admin" {
		gonic.Gonic(cherry.ErrAdminRequired(), ctx)
		return
	}
}
