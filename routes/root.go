package routes

import (
	"net/http"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"github.com/gin-gonic/gin"
)

func simpleSendHandler(ctx *gin.Context) {
	var request mttypes.SimpleSendRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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
	info, err := svc.UserManagerClient.UserInfoByID(ctx, request.UserID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	recipient := &mttypes.Recipient{
		ID:    request.UserID,
		Name:  info.Login,
		Email: info.Login,
	}
	status, err := svc.UpstreamSimple.SimpleSend(ctx, request.Template, tv, recipient)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, mttypes.SimpleSendResponse{
		UserID: status.RecipientID,
	})
}
