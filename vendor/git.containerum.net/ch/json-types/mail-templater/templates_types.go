package mail

import "time"

//Types related to storing templates

type Template struct {
	Name      string    `json:"name,omitempty"`
	Version   string    `json:"version,omitempty"`
	Data      string    `json:"data,omitempty"`
	Subject   string    `json:"subject,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type TemplatesListEntry struct {
	Name     string   `json:"name,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

type TemplatesListResponse struct {
	Templates []TemplatesListEntry `json:"templates,omitempty"`
}