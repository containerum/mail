package httputil

import (
	"github.com/gin-gonic/gin"

)

type Masker interface {
	Mask()
}

func MaskForNonAdmin(ctx *gin.Context, m Masker) {
	if ctx.GetHeader(UserRoleXHeader) != "admin" {
		m.Mask()
	}
}
