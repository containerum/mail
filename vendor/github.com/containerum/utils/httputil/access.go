package httputil

import (
	"context"

	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	kubeModel "github.com/containerum/kube-client/pkg/model"
	"github.com/gin-gonic/gin"
)

const AccessContext = "access-ctx"
const AllAccessContext = "all-access-ctx"

const (
	ProjectParam   = "project"
	NamespaceParam = "namespace"
)

type ProjectAccess struct {
	ProjectID          string            `json:"project_id"`
	ProjectLabel       string            `json:"project_label"`
	NamespacesAccesses []NamespaceAccess `json:"namespaces"`
}

type NamespaceAccess struct {
	NamespaceID    string                    `json:"namespace_id"`
	NamespaceLabel string                    `json:"namespace_label"`
	Access         kubeModel.UserGroupAccess `json:"access"`
}

type Permissions interface {
	GetAllAccesses(ctx context.Context) ([]ProjectAccess, error)
	GetNamespaceAccess(ctx context.Context, projectID, namespaceID string) (NamespaceAccess, error)
}

type AccessChecker struct {
	PermissionsClient Permissions
	AccessError       cherry.ErrConstruct
	NotFoundError     cherry.ErrConstruct
}

func (a *AccessChecker) CheckAccess(requiredAccess kubeModel.UserGroupAccess) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if MustGetUserRole(ctx.Request.Context()) == "admin" {
			return
		}
		project := ctx.Param(ProjectParam)
		ns := ctx.Param(NamespaceParam)

		namespaceAccess, err := a.PermissionsClient.GetNamespaceAccess(ctx.Request.Context(), project, ns)
		if err != nil {
			gonic.Gonic(a.AccessError(), ctx)
			return
		}

		if namespaceAccess.Access < requiredAccess {
			gonic.Gonic(a.NotFoundError(), ctx)
			return
		}

		rctx := context.WithValue(ctx.Request.Context(), AccessContext, namespaceAccess)
		ctx.Request = ctx.Request.WithContext(rctx)
	}
}

func (a *AccessChecker) SaveAllAccesses() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		namespaceAccess, err := a.PermissionsClient.GetAllAccesses(ctx.Request.Context())
		if err != nil {
			gonic.Gonic(a.AccessError(), ctx)
			return
		}

		rctx := context.WithValue(ctx.Request.Context(), AllAccessContext, namespaceAccess)
		ctx.Request = ctx.Request.WithContext(rctx)
	}
}
