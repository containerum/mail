package models

// PasswordChangeRequest -- password change request
//
// swagger:model
type PasswordChangeRequest struct {
	// required: true
	CurrentPassword string `json:"current_password"`
	// required: true
	NewPassword string `json:"new_password"`
}

// PasswordRestoreRequest -- password restore request
//
// swagger:model
type PasswordRestoreRequest struct {
	// required: true
	Link string `json:"link"`
	// required: true
	NewPassword string `json:"new_password"`
}
