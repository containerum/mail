package routes

import (
	"github.com/gin-gonic/gin"
)

type ReCaptchaRequest struct {
	ReCaptcha string `json:"recaptcha" binding:"required"`
}

const (
	reCaptchaFailed = "reCaptcha failed"
)

func reCaptchaMiddleware(ctx *gin.Context) {
	/*var request ReCaptchaRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	checkResp, err := svc.ReCaptchaClient.Check(ctx.ClientIP(), request.ReCaptcha)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !checkResp.Success {
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewError(reCaptchaFailed))
	}*/
}
