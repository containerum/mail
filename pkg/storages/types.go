package storages

import (
	"io"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
)

// TemplateStorage used to store email templates.
// Implementation must support tagging and versioning with "semantic versioning 2.0"
type TemplateStorage interface {
	// PutTemplate puts template to storage.
	PutTemplate(templateName, templateVersion, templateData, templateSubject string, new bool) error

	// GetTemplate returns specified version of template.
	GetTemplate(templateName, templateVersion string) (*mttypes.Template, error)

	// GetLatestVersionTemplate returns latest version of template and it`s value using semver to compare versions.
	GetLatestVersionTemplate(templateName string) (*string, *mttypes.Template, error)

	// GetTemplates returns all versions of templates in map (key is version, value is template).
	GetTemplates(templateName string) (map[string]*mttypes.Template, error)

	// DeleteTemplate deletes specified version of template. Returns nil on successful delete.
	DeleteTemplate(templateName, templateVersion string) error

	// DeleteTemplates deletes all versions of template. Returns nil on successful delete.
	DeleteTemplates(templateName string) error

	// GetTemplatesList returns list of all of templates.
	GetTemplatesList() (*mttypes.TemplatesListResponse, error)

	io.Closer
}

// MessagesStorage used to store sent emails.
type MessagesStorage interface {
	// PutMessage puts MessageStorageValue to storage.
	// If message with specified id already exists in storage it will be overwritten.
	PutMessage(id string, value *mttypes.MessagesStorageValue) error

	// GetMessage returns value by specified ID.
	GetMessage(id string) (*mttypes.MessagesStorageValue, error)

	// GetValue returns all messages.
	GetMessageList(page int, perPage int) (*mttypes.MessageListResponse, error)

	io.Closer
}
