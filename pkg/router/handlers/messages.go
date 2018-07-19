package handlers

import (
	"net/http"

	"github.com/containerum/cherry"
	"github.com/gin-gonic/gin"

	"strconv"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"github.com/containerum/cherry/adaptors/gonic"
)

// swagger:operation GET /messages Messages MessageListGetHandler
//// Get messages list.
//// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#get-messages
////
//// ---
//// x-method-visibility: private
//// parameters:
////  - $ref: '#/parameters/UserRoleHeader'
////  - name: page
////    in: query
////    type: string
////    required: false
////  - name: per_page
////    in: query
////    type: string
////    required: false
//// responses:
////  '200':
////    description: message list get response
////    schema:
////      $ref: '#/definitions/MessageListResponse'
////  default:
////    $ref: '#/responses/error'
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

// swagger:operation GET /messages/{message_id} Messages MessageGetHandler
// Get single message.
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#get-message-copy
//
// ---
// x-method-visibility: private
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - name: message_id
//    in: path
//    type: string
//    required: true
// responses:
//  '200':
//    description: message get response
//    schema:
//      $ref: '#/definitions/MessageGetResponse'
//  default:
//    $ref: '#/responses/error'
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
		ID:                   id,
		MessagesStorageValue: v,
	})
}
