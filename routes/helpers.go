package routes

import (
	"net/http"

	"git.containerum.net/ch/mail-templater/storages"
	chutils "git.containerum.net/ch/utils"
	"github.com/gin-gonic/gin"
)

func sendStorageError(ctx *gin.Context, err error) {
	switch err {
	case nil:
	case storages.ErrTemplateNotExists, storages.ErrVersionNotExists, storages.ErrMessageNotExists:
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewError(err.Error()))
	default:
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func sendValidationError(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
}
