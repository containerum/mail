package routes

import (
	"net/http"

	"time"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/grpc-proto-files/common"
	"git.containerum.net/ch/mail-templater/upstreams"
	"git.containerum.net/ch/user-manager/models"
	"git.containerum.net/ch/user-manager/utils"
	chutils "git.containerum.net/ch/utils"
	"github.com/gin-gonic/gin"
)

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type PasswordRestoreRequest struct {
	Link        string `json:"link" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

const (
	invalidPassword    = "invalid password provided"
	linkNotForPassword = "link %s is not for password changing"
	userBanned         = "user %s banned"
)

func passwordChangeHandler(ctx *gin.Context) {
	userID := ctx.GetHeader(UserIDHeader)
	var request PasswordChangeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	user, err := svc.DB.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userWithIDNotFound, userID))
		return
	}

	if !utils.CheckPassword(request.CurrentPassword, user.Salt, user.PasswordHash) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewError(invalidPassword))
		return
	}

	_, err = svc.AuthClient.DeleteUserTokens(ctx, &auth.DeleteUserTokensRequest{
		UserId: &common.UUID{Value: user.ID},
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user.PasswordHash = utils.GetKey(request.NewPassword, user.Salt)
	err = svc.DB.UpdateUser(user)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = svc.MailClient.SendPasswordChangedMail(&upstreams.Recipient{
		ID:        user.ID,
		Name:      user.Login,
		Email:     user.Login,
		Variables: map[string]string{},
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO: get access from resource manager

	tokens, err := svc.AuthClient.CreateToken(ctx, &auth.CreateTokenRequest{
		UserAgent:   ctx.Request.UserAgent(),
		UserId:      &common.UUID{Value: user.ID},
		UserIp:      ctx.ClientIP(),
		UserRole:    auth.Role(user.Role),
		RwAccess:    true,
		Access:      &auth.ResourcesAccess{},
		PartTokenId: nil,
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusAccepted, tokens)
}

func passwordResetHandler(ctx *gin.Context) {
	userID := ctx.GetHeader(UserIDHeader)
	var request PasswordChangeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	user, err := svc.DB.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userWithIDNotFound, userID))
		return
	}
	if user.IsInBlacklist {
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewErrorF(userBanned, user.Login))
		return
	}

	var link *models.Link
	err = svc.DB.Transactional(func(tx *models.DB) (err error) {
		link, err = svc.DB.CreateLink(models.LinkTypePwdChange, 24*time.Hour, user)
		return
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = svc.MailClient.SendPasswordResetMail(&upstreams.Recipient{
		ID:        user.ID,
		Name:      user.Login,
		Email:     user.Login,
		Variables: map[string]string{"TOKEN": link.Link},
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func passwordRestoreHandler(ctx *gin.Context) {
	var request PasswordRestoreRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	link, err := svc.DB.GetLinkFromString(request.Link)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if link == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(linkNotFound, request.Link))
		return
	}
	if link.Type != models.LinkTypePwdChange {
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewErrorF(linkNotForPassword, request.Link))
		return
	}

	_, err = svc.AuthClient.DeleteUserTokens(ctx, &auth.DeleteUserTokensRequest{
		UserId: &common.UUID{Value: link.User.ID},
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	link.User.PasswordHash = utils.GetKey(request.NewPassword, link.User.Salt)
	err = svc.DB.UpdateUser(link.User)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = svc.MailClient.SendPasswordChangedMail(&upstreams.Recipient{
		ID:        link.User.ID,
		Name:      link.User.Login,
		Email:     link.User.Login,
		Variables: map[string]string{},
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO: get access from resource manager

	tokens, err := svc.AuthClient.CreateToken(ctx, &auth.CreateTokenRequest{
		UserAgent:   ctx.Request.UserAgent(),
		UserId:      &common.UUID{Value: link.User.ID},
		UserIp:      ctx.ClientIP(),
		UserRole:    auth.Role(link.User.Role),
		RwAccess:    true,
		Access:      &auth.ResourcesAccess{},
		PartTokenId: nil,
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusAccepted, tokens)
}
