package routes

import (
	"net/http"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"github.com/gin-gonic/gin"
)

type SimpleSendRequest struct {
	Template  string            `json:"template" binding:"required"`
	UserID    string            `json:"user_id" binding:"required,uuid4"`
	Variables map[string]string `json:"variables"`
}

type SimpleSendResponse struct {
	UserID string `json:"user_id"`
}

func simpleSendHandler(ctx *gin.Context) {
	var request SimpleSendRequest
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
	info, err := svc.UserManagerClient.UserInfoByID(request.UserID)
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
	status, err := svc.Upstream.SimpleSend(request.Template, tv, recipient)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, &SimpleSendResponse{
		UserID: status.RecipientID,
	})
	ctx.JSON(http.StatusOK, &SimpleSendResponse{
		UserID: request.UserID,
	})
}
