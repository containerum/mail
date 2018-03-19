package gonic

import (
	"fmt"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"github.com/gin-gonic/gin"
)

// Gonic -- aborts gin HTTP request with StatusHTTP
// and provides json representation of error
func Gonic(err *cherry.Err, ctx *gin.Context) {
	ctx.Error(err)
	ctx.AbortWithStatusJSON(err.StatusHTTP, err)
}

// Recovery -- gin middleware to catch panics and wrap it to cherry error.
// If panic caught it aborts HTTP request with defaultErr.
func Recovery(defaultErr cherry.ErrConstruct, logger cherry.ErrorLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			var errToReturn *cherry.Err
			if r := recover(); r != nil {
				if cherryErr, ok := r.(*cherry.Err); ok {
					errToReturn = cherryErr
				} else {
					errToReturn = defaultErr().Log(fmt.Errorf("%v", r), logger)
				}
				Gonic(errToReturn, ctx)
			}
		}()

		ctx.Next()
	}
}
