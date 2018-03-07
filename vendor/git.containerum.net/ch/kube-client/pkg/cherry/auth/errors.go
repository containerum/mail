package autherr

import (
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/cherry"
)

var buildErr = cherry.BuildErr(cherry.Auth)

func ErrInvalidToken() *cherry.Err {
	return buildErr("invalid token received", http.StatusBadRequest, 1)
}

func ErrTokenNotOwnedBySender() *cherry.Err {
	return buildErr("can`t identify sender as token owner", http.StatusForbidden, 2)
}

func ErrTokenNotFound() *cherry.Err {
	return buildErr("token was not found in storage", http.StatusNotFound, 3)
}

func ErrInternal() *cherry.Err {
	return buildErr("internal error", http.StatusInternalServerError, 4)
}

func ErrValidation() *cherry.Err {
	return buildErr("validation error", http.StatusBadRequest, 5)
}
