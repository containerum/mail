package routes

import (
	"net/http"

	"bitbucket.org/exonch/ch-mail-templater/upstreams"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type simpleSendRequest struct {
	Template  string            `json:"template" binding:"required"`
	UserID    string            `json:"user_id" binding:"required"`
	Variables map[string]string `json:"variables" binding:"required"`
}

type simpleSendResponse struct {
	UserID string `json:"user_id"`
}

func simpleSendHandler(ctx *gin.Context) {
	var request simpleSendRequest
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.Error(err)
		sendValidationError(ctx, err)
		return
	}
	tv, err := svc.TemplateStorage.GetTemplate(request.Template, "latest") // TODO: properly handle
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	// TODO: here we should contact with user manager and get email by id
	recipient := &upstreams.Recipient{}
	status, err := svc.Upstream.SimpleSend(request.Template, tv, recipient)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, &simpleSendResponse{
		UserID: status.RecipientID,
	})
}
