package routes

import (
	"net/http"

	"git.containerum.net/ch/json-types/errors"
	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"github.com/gin-gonic/gin"
)

func sendStorageError(ctx *gin.Context, err error) {
	switch err {
	case nil:
	case mttypes.ErrTemplateNotExists, mttypes.ErrVersionNotExists, mttypes.ErrMessageNotExists:
		ctx.AbortWithStatusJSON(http.StatusNotFound, []*errors.Error{errors.New(err.Error())})
	default:
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
