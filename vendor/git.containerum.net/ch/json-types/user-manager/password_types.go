package user

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type PasswordRestoreRequest struct {
	Link        string `json:"link" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type PasswordResetRequest struct {
	Username string `json:"username" binding:"required,email"`
}