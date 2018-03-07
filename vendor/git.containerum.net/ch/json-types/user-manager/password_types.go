package user

type PasswordRequest struct {
	Link            string `json:"link"`
	Token           string `json:"token"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
