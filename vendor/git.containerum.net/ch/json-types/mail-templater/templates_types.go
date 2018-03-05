package mail

import "time"

//Types related to storing templates

type Template struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Data      string    `json:"data"`
	Subject   string    `json:"ubject"`
	CreatedAt time.Time `json:"created_at"`
}

type TemplatesListEntry struct {
	Name     string   `json:"template_name"`
	Versions []string `json:"template_versions"`
}

type TemplatesListResponse struct {
	Templates []TemplatesListEntry `json:"templates"`
}