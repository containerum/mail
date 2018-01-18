package routes

import (
	"net/http"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"github.com/gin-gonic/gin"
)

func messageGetHandler(ctx *gin.Context) {
	id := ctx.Param("message_id")
	v, err := svc.MessagesStorage.GetValue(id)
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &mttypes.MessageGetResponse{
		Id:                   id,
		MessagesStorageValue: v,
	})
}
