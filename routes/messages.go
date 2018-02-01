package routes

import (
	"net/http"

	"git.containerum.net/ch/json-types/errors"
	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"github.com/gin-gonic/gin"
)

func messageGetHandler(ctx *gin.Context) {
	id := ctx.Param("message_id")
	v, err := svc.MessagesStorage.GetValue(id)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
		return
	}
	ctx.JSON(http.StatusOK, &mttypes.MessageGetResponse{
		Id:                   id,
		MessagesStorageValue: v,
	})
}

func messageListGetHandler(ctx *gin.Context) {
	var params mttypes.MessageListQuery
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ParseBindErorrs(err))
		return
	}

	v, err := svc.MessagesStorage.GetMessageList(params.Page, params.PerPage)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
		return
	}
	ctx.JSON(http.StatusOK, v)
}
