package handlers

import (
	"net/http"

	"strconv"

	"git.containerum.net/ch/cherry"
	"github.com/gin-gonic/gin"

	"git.containerum.net/ch/cherry/adaptors/gonic"
	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
)

//MessageGetHandler returns one message
func MessageGetHandler(ctx *gin.Context) {
	id := ctx.Param("message_id")

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	v, err := svc.MessagesStorage.GetMessage(id)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrUnableGetMessage(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, &models.MessageGetResponse{
		Id:                   id,
		MessagesStorageValue: v,
	})
}

//MessageListGetHandler returns messages list
func MessageListGetHandler(ctx *gin.Context) {
	page := int64(1)
	pagestr, ok := ctx.GetQuery("page")
	if ok {
		var err error
		page, err = strconv.ParseInt(pagestr, 10, 64)
		if err != nil {
			ctx.Error(err)
		}
	}

	perPage := int64(10)
	perPagestr, ok := ctx.GetQuery("per_page")
	if ok {
		var err error
		perPage, err = strconv.ParseInt(perPagestr, 10, 64)
		if err != nil {
			ctx.Error(err)
		}
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	v, err := svc.MessagesStorage.GetMessageList(int(page), int(perPage))
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrUnableGetMessagesList(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, v)
}
