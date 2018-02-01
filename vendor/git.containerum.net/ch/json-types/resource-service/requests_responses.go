package resource

import (
	"reflect"

	"regexp"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

type CreateResourceRequest struct {
	TariffID string `json:"tariff-id" binding:"uuid4"`
	Label    string `json:"label" binding:"required,dns"`
}

type RenameResourceRequest struct {
	NewLabel string `json:"label" binding:"required,dns"`
}

type GetAllResourcesQueryParams struct {
	Page    int    `form:"page" binding:"gt=0"`
	PerPage int    `form:"per_page" binding:"gt=0"`
	Filters string `form:"filters"`
}

type SetResourcesAccessRequest struct {
	Access PermissionStatus `json:"access" binding:"eq=owner|eq=read|eq=write|eq=readdelete|eq=none"`
}

type ResizeResourceRequest struct {
	NewTariffID string `json:"tariff_id" binding:"uuid4"`
}

// Namespaces

type CreateNamespaceRequest = CreateResourceRequest

type GetUserNamespacesResponse []NamespaceWithVolumes

func (r GetUserNamespacesResponse) Mask() {
	for i := range r {
		r[i].Mask()
	}
}

type GetUserNamespaceResponse = NamespaceWithVolumes

type GetAllNamespacesResponse []NamespaceWithVolumes

func (r GetAllNamespacesResponse) Mask() {
	for i := range r {
		r[i].Mask()
	}
}

type GetUserNamespaceAccessesResponse = NamespaceWithUserPermissions

type RenameNamespaceRequest = RenameResourceRequest

type SetNamespaceAccessRequest struct {
	Username string           `json:"username" binding:"email"`
	Access   PermissionStatus `json:"access" binding:"eq=owner|eq=read|eq=write|eq=readdelete|eq=none"`
}

type ResizeNamespaceRequest = ResizeResourceRequest

// Volumes

type CreateVolumeRequest = CreateResourceRequest

type GetUserVolumesResponse []VolumeWithPermission

func (r GetUserVolumesResponse) Mask() {
	for i := range r {
		r[i].Mask()
	}
}

type GetUserVolumeResponse = VolumeWithPermission

type GetAllVolumesResponse []VolumeWithPermission

func (r GetAllVolumesResponse) Mask() {
	for i := range r {
		r[i].Mask()
	}
}

type GetVolumeAccessesResponse = VolumeWithUserPermissions

type RenameVolumeRequest = RenameResourceRequest

type SetVolumeAccessRequest = SetNamespaceAccessRequest

type ResizeVolumeRequest = ResizeResourceRequest

// Other

// GetUserAccessResponse is response for special request needed for auth server (actually for creating tokens)
type GetUserAccessesResponse = auth.ResourcesAccess

// custom tag registration

func RegisterCustomTags(validate *validator.Validate) error {
	return validate.RegisterValidation("dns", dnsValidationFunc)
}

func RegisterCustomTagsGin(validate binding.StructValidator) error {
	return validate.RegisterValidation("dns", dnsValidationFunc)
}

var dnsLabel = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

func dnsValidationFunc(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	return dnsLabel.MatchString(field.String())
}
