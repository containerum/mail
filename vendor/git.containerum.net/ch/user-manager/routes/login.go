package routes

import (
	"net/http"

	"time"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/grpc-proto-files/common"
	"git.containerum.net/ch/mail-templater/upstreams"
	"git.containerum.net/ch/user-manager/clients"
	"git.containerum.net/ch/user-manager/models"
	chutils "git.containerum.net/ch/utils"
	"github.com/gin-gonic/gin"
)

type BasicLoginRequest struct {
	Login     string `json:"login" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	ReCaptcha string `json:"recaptcha" binding:"required"`
}

type OneTimeTokenLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

type OAuthLoginRequest struct {
	Resource    clients.OAuthResource `json:"resource" binding:"required"`
	AccessToken string                `json:"access_token" binding:"required"`
}

const (
	oneTimeTokenNotFound = "one-time token %s not exists or already used"
	resourceNotSupported = "resource %s not supported now"
)

func basicLoginHandler(ctx *gin.Context) {
	var request BasicLoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	user, err := svc.DB.GetUserByLogin(request.Login)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userNotFound, request.Login))
		return
	}
	if user.IsInBlacklist {
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewErrorF(userBanned, request.Login))
		return
	}

	if !user.IsActive {
		link, err := svc.DB.GetLinkForUser(models.LinkTypeConfirm, user)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if link == nil {
			link, err = svc.DB.CreateLink(models.LinkTypeConfirm, 24*time.Hour, user)
			if err != nil {
				ctx.Error(err)
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		if tdiff := time.Now().UTC().Sub(link.SentAt.Time); link.SentAt.Valid && tdiff < 5*time.Minute {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewErrorF(waitForResend, int(tdiff.Seconds())))
			return
		}

		err = svc.MailClient.SendConfirmationMail(&upstreams.Recipient{
			ID:        user.ID,
			Name:      user.Login,
			Email:     user.Login,
			Variables: map[string]string{"CONFIRM": link.Link},
		})
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
		return
	}

	// TODO: get access data from resource manager
	access := &auth.ResourcesAccess{}

	tokens, err := svc.AuthClient.CreateToken(ctx, &auth.CreateTokenRequest{
		UserAgent:   ctx.Request.UserAgent(),
		UserId:      &common.UUID{Value: user.ID},
		UserIp:      ctx.ClientIP(),
		UserRole:    auth.Role(user.Role),
		RwAccess:    true,
		Access:      access,
		PartTokenId: nil,
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func oneTimeTokenLoginHandler(ctx *gin.Context) {
	var request OneTimeTokenLoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	token, err := svc.DB.GetTokenObject(request.Token)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if token == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(oneTimeTokenNotFound, request.Token))
		return
	}
	if token.User.IsInBlacklist {
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewErrorF(userBanned, token.User.Login))
		return
	}

	// TODO: get access data from resource manager
	access := &auth.ResourcesAccess{}

	var tokens *auth.CreateTokenResponse

	err = svc.DB.Transactional(func(tx *models.DB) error {
		var err error
		tokens, err = svc.AuthClient.CreateToken(ctx, &auth.CreateTokenRequest{
			UserAgent:   ctx.Request.UserAgent(),
			UserId:      &common.UUID{Value: token.User.ID},
			UserIp:      ctx.ClientIP(),
			UserRole:    auth.Role(token.User.Role),
			RwAccess:    true,
			Access:      access,
			PartTokenId: nil,
		})
		if err != nil {
			return err
		}

		token.IsActive = false
		token.SessionID = "sid" // TODO: session ID here
		if err := tx.UpdateToken(token); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func oauthLoginHandler(ctx *gin.Context) {
	var request OAuthLoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	resource, exist := clients.OAuthClientByResource(request.Resource)
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewErrorF(resourceNotSupported, request.Resource))
		return
	}

	info, err := resource.GetUserInfo(request.AccessToken)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := svc.DB.GetUserByLogin(info.Email)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userNotFound, info.Email))
		return
	}
	if user.IsInBlacklist {
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewErrorF(userBanned, user.Login))
		return
	}

	accounts, err := svc.DB.GetUserBoundAccounts(user)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if accounts == nil {
		if err := svc.DB.Transactional(func(tx *models.DB) error {
			return tx.BindAccount(user, string(request.Resource), info.UserID)
		}); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	// TODO: get access data from resource manager
	access := &auth.ResourcesAccess{}

	tokens, err := svc.AuthClient.CreateToken(ctx, &auth.CreateTokenRequest{
		UserAgent:   ctx.Request.UserAgent(),
		UserId:      &common.UUID{Value: user.ID},
		UserIp:      ctx.ClientIP(),
		UserRole:    auth.Role(user.Role),
		RwAccess:    true,
		Access:      access,
		PartTokenId: nil,
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func webAPILoginHandler(ctx *gin.Context) {
	var request clients.WebAPILoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	resp, code, err := svc.WebAPIClient.Login(&request)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(code, chutils.NewError(err.Error()))
		return
	}

	// TODO: get access data from resource manager
	access := &auth.ResourcesAccess{}

	tokens, err := svc.AuthClient.CreateToken(ctx, &auth.CreateTokenRequest{
		UserAgent:   ctx.Request.UserAgent(),
		UserId:      &common.UUID{Value: resp["user"].(map[string]interface{})["id"].(string)},
		UserIp:      ctx.ClientIP(),
		UserRole:    auth.Role_USER,
		RwAccess:    true,
		Access:      access,
		PartTokenId: nil,
	})

	resp["access_token"] = tokens.AccessToken
	resp["refresh_token"] = tokens.RefreshToken

	ctx.JSON(http.StatusOK, resp)
}
