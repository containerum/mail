package storages

import (
	"io"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
)

// TemplateStorage used to store email templates.
// Implementation must support tagging and versioning with "semantic versioning 2.0"
type TemplateStorage interface {
	// PutTemplate puts template to storage. If template with specified name and version exists it will be overwritten.
	PutTemplate(templateName, templateVersion, templateData, templateSubject string) error

	// GetTemplate returns specified version of template.
	GetTemplate(templateName, templateVersion string) (*mttypes.TemplateStorageValue, error)

	// GetLatestVersionTemplate returns latest version of template and it`s value using semver to compare versions.
	GetLatestVersionTemplate(templateName string) (string, *mttypes.TemplateStorageValue, error)

	// GetTemplates returns all versions of templates in map (key is version, value is template).
	GetTemplates(templateName string) (map[string]*mttypes.TemplateStorageValue, error)

	// DeleteTemplate deletes specified version of template. Returns nil on successful delete.
	DeleteTemplate(templateName, templateVersion string) error

	// DeleteTemplates deletes all versions of template. Returns nil on successful delete.
	DeleteTemplates(templateName string) error

	io.Closer
}

// MessagesStorage used to store sent emails.
type MessagesStorage interface {
	// PutValue puts MessageStorageValue to storage.
	// If message with specified id already exists in storage it will be overwritten.
	PutValue(id string, value *mttypes.MessagesStorageValue) error

	// GetValue returns value by specified ID.
	GetValue(id string) (*mttypes.MessagesStorageValue, error)

	io.Closer
}
