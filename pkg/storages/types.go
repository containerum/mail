package storages

import (
	"io"

	"git.containerum.net/ch/mail-templater/pkg/models"
)

// TemplateStorage used to store email templates.
// Implementation must support tagging and versioning with "semantic versioning 2.0"
type TemplateStorage interface {
	// PutTemplate puts template to storage.
	PutTemplate(templateName, templateVersion, templateData, templateSubject string, new bool) error

	// GetTemplate returns specified version of template.
	GetTemplate(templateName, templateVersion string) (*models.Template, error)

	// GetLatestVersionTemplate returns latest version of template and it`s value using semver to compare versions.
	GetLatestVersionTemplate(templateName string) (*string, *models.Template, error)

	// GetTemplates returns all versions of templates in map (key is version, value is template).
	GetTemplates(templateName string) (map[string]*models.Template, error)

	// DeleteTemplate deletes specified version of template. Returns nil on successful delete.
	DeleteTemplate(templateName, templateVersion string) error

	// DeleteTemplates deletes all versions of template. Returns nil on successful delete.
	DeleteTemplates(templateName string) error

	// GetTemplatesList returns list of all of templates.
	GetTemplatesList() (*models.TemplatesListResponse, error)

	io.Closer
}

// MessagesStorage used to store sent emails.
type MessagesStorage interface {
	// PutMessage puts MessageStorageValue to storage.
	// If message with specified id already exists in storage it will be overwritten.
	PutMessage(id string, value *models.MessagesStorageValue) error

	// GetMessage returns value by specified ID.
	GetMessage(id string) (*models.MessagesStorageValue, error)

	// GetValue returns all messages.
	GetMessageList(page int, perPage int) (*models.MessageListResponse, error)

	io.Closer
}
