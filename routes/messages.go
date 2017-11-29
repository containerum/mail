package routes

import (
	"net/http"

	"git.containerum.net/ch/mail-templater/storages"
	"github.com/gin-gonic/gin"
)

type messageGetResponse struct {
	Id string `json:"id"`
	*storages.MessagesStorageValue
}

func messageGetHandler(ctx *gin.Context) {
	id := ctx.Param("message_id")
	v, err := svc.MessagesStorage.GetValue(id)
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, messageGetResponse{
		Id:                   id,
		MessagesStorageValue: v,
	})
}
