package routes

import (
	"net/http"

	"git.containerum.net/ch/json-types/errors"
	umtypes "git.containerum.net/ch/json-types/user-manager"
	"github.com/gin-gonic/gin"
)

func requireHeaders(headers ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var notFoundHeaders []string
		for _, v := range headers {
			if ctx.GetHeader(v) == "" {
				notFoundHeaders = append(notFoundHeaders, v)
			}
		}
		if len(notFoundHeaders) > 0 {
			err := errors.Format("required headers %v was not provided", notFoundHeaders)
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
	}
}

func requireAdminRole(ctx *gin.Context) {
	if ctx.GetHeader(umtypes.UserRoleHeader) != "admin" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, errors.New("Only admin can do this"))
		return
	}

	userID := ctx.GetHeader(umtypes.UserIDHeader)

	info, err := svc.UserManagerClient.UserInfoByID(ctx, userID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, err.Error())
		return
	}

	if info != nil {
		if info.Role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.New("Only admin can do this"))
			return
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusForbidden, errors.New("Unable to verify your permissions"))
	}
}
