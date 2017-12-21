package routes

import (
	"net/http"

	"git.containerum.net/ch/auth/storages"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/grpc-proto-files/common"
	chutils "git.containerum.net/ch/utils"
	"github.com/gin-gonic/gin"
)

const (
	tokenNotOwnedByUser = "token %s not owned by user %s"
)

func logoutHandler(ctx *gin.Context) {
	tokenID := ctx.Param("token_id")
	userID := ctx.GetHeader(UserIDHeader)
	_, err := svc.AuthClient.DeleteToken(ctx, &auth.DeleteTokenRequest{
		TokenId: &common.UUID{Value: tokenID},
		UserId:  &common.UUID{Value: userID},
	})

	switch err {
	case nil:
	case storages.ErrTokenNotOwnedBySender:
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewError(err.Error()))
		return
	case storages.ErrInvalidToken:
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	default:
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	oneTimeToken, err := svc.DB.GetTokenBySessionID(ctx.GetHeader(SessionIDHeader))
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if oneTimeToken != nil {
		if oneTimeToken.User.ID != userID {
			ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewErrorF(tokenNotOwnedByUser, oneTimeToken.Token, userID))
			return
		}
		if err := svc.DB.DeleteToken(oneTimeToken.Token); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	ctx.Status(http.StatusOK)
}
