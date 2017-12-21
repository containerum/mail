package routes

import (
	"net/http"

	"git.containerum.net/ch/mail-templater/upstreams"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type simpleSendRequest struct {
	Template  string            `json:"template" binding:"required"`
	UserID    string            `json:"user_id" binding:"required,uuid4"`
	Variables map[string]string `json:"variables"`
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
	_, tv, err := svc.TemplateStorage.GetLatestVersionTemplate(request.Template)
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	info, err := svc.UserManagerClient.UserInfoByID(request.UserID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	recipient := &upstreams.Recipient{
		ID:    request.UserID,
		Name:  info.Login,
		Email: info.Login,
	}
	status, err := svc.Upstream.SimpleSend(request.Template, tv, recipient)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, &simpleSendResponse{
		UserID: status.RecipientID,
	})
	ctx.JSON(http.StatusOK, &simpleSendResponse{
		UserID: request.UserID,
	})
}
