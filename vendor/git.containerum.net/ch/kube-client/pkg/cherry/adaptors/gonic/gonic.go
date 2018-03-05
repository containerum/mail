package gonic

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"github.com/gin-gonic/gin"
)

// Gonic -- aborts gin HTTP request with StatusHTTP
// and provides json representation of error
func Gonic(err *cherry.Err, ctx *gin.Context) {
	ctx.Error(err)
	ctx.AbortWithStatusJSON(err.StatusHTTP, err)
}
