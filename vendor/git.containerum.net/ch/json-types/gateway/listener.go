package gateway

type ListenerJSON struct {
	ID        *string `json:"id,omitempty"`
	CreatedAt *int64  `json:"created_at,omitempty"`
	UpdatedAt *int64  `json:"updated_at,omitempty"`
	DeletedAt *int64  `json:"deleted_at,omitempty"`
	Name      *string `json:"name,omitempty"`
	// Roles       []Role   `json:"deleted_at,omitempty"`
	OAuth       *bool      `json:"o_auth,omitempty"`
	Active      *bool      `json:"active,omitempty"`
	Group       *GroupJSON `json:"group,omitempty"`
	GroupID     *string    `json:"group_id,omitempty"`
	StripPath   *bool      `json:"strip_path,omitempty"`
	ListenPath  *string    `json:"listen_path,omitempty"`
	UpstreamURL *string    `json:"upstream_url,omitempty"`
	Method      *string    `json:"method,omitempty"`
	// Plugins     []Plugin  `json:"deleted_at,omitempty"`
}
