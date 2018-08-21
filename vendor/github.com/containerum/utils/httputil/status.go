package httputil

import (
	"net/http"

	kubeModel "github.com/containerum/kube-client/pkg/model"
	"github.com/gin-gonic/gin"
)

func ServiceStatus(status *kubeModel.ServiceStatus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var httpStatus int
		if status.StatusOK {
			httpStatus = http.StatusOK
		} else {
			httpStatus = http.StatusInternalServerError
		}
		ctx.JSON(httpStatus, status)
	}
}
