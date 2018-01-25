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

type SetResourceAccessRequest struct {
	UserID string           `json:"user_id" binding:"uuid4"`
	Access PermissionStatus `json:"access"`
}

type ResizeResourceRequest struct {
	NewTariffID string `json:"tariff_id" binding:"uuid4"`
}

// Namespaces

type CreateNamespaceRequest = CreateResourceRequest

type GetUserNamespacesResponse = []NamespaceWithVolumes

type GetUserNamespaceResponse = NamespaceWithVolumes

type GetAllNamespacesResponse = []NamespaceWithVolumes

type GetNamespaceAccessResponse = NamespaceWithPermission

type RenameNamespaceRequest = RenameResourceRequest

type SetNamespaceAccessRequest = SetResourceAccessRequest

type ResizeNamespaceRequest = ResizeResourceRequest

// Volumes

type CreateVolumeRequest = CreateResourceRequest

type GetUserVolumesRequest = []VolumeWithPermission

type GetUserVolumeRequest = VolumeWithPermission

type GetAllVolumesRequest = []VolumeWithPermission

type GetVolumeAccessesResponse = VolumeWithUserPermissions

type RenameVolumeRequest = RenameResourceRequest

type SetVolumeAccessRequest = SetResourceAccessRequest

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
