package models

import "time"

// Template -- template model
//
// swagger:model
type Template struct {
	// required: true
	Name string `json:"name,omitempty"`
	// required: true
	Version string `json:"version,omitempty"`
	// required: true
	Data string `json:"data,omitempty"`
	// required: true
	Subject   string     `json:"subject,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// TemplatesListEntry -- model for template list
//
// swagger:model
type TemplatesListEntry struct {
	Name     string   `json:"name,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

// TemplatesListResponse -- templates list response
//
// swagger:model
type TemplatesListResponse struct {
	Templates []TemplatesListEntry `json:"templates,omitempty"`
}
