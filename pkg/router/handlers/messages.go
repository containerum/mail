package handlers

import (
	"net/http"

	"strconv"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	ch "git.containerum.net/ch/kube-client/pkg/cherry"
	"github.com/gin-gonic/gin"

	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/mail-templater"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
)

func MessageGetHandler(ctx *gin.Context) {
	id := ctx.Param("message_id")

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	v, err := svc.MessagesStorage.GetMessage(id)
	if err != nil {
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableGetMessage(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, &mttypes.MessageGetResponse{
		Id:                   id,
		MessagesStorageValue: v,
	})
}

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
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableGetMessagesList(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, v)
}
