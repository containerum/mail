package routes

import (
	"net/http"

	"git.containerum.net/ch/mail-templater/storages"
	"github.com/gin-gonic/gin"
)

type errorToClient struct {
	Error string `json:"error"`
}

func sendStorageError(ctx *gin.Context, err error) {
	switch err {
	case nil:
	case storages.ErrTemplateNotExists, storages.ErrVersionNotExists, storages.ErrMessageNotExists:
		ctx.AbortWithStatusJSON(http.StatusNotFound, errorToClient{Error: err.Error()})
	default:
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func sendValidationError(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, errorToClient{Error: err.Error()})
}
