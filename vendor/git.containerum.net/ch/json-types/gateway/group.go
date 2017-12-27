package gateway

//GroupJSON middleware struct for Group Marshal and Unmarshal
type GroupJSON struct {
	ID        *string `json:"id,omitempty"`
	CreatedAt *int64  `json:"created_at,omitempty"`
	UpdatedAt *int64  `json:"updated_at,omitempty"`
	DeletedAt *int64  `json:"deleted_at,omitempty"`
	Name      *string `json:"name,omitempty"`
	Active    *bool   `json:"active,omitempty"`
}
