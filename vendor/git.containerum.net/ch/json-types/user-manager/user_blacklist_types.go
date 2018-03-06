package user

type BlacklistedUserEntry struct {
	Login string `json:"login"`
	ID    string `json:"id"`
}

type BlacklistGetResponse struct {
	BlacklistedUsers []BlacklistedUserEntry `json:"blacklist_users"`
}

type UserToBlacklistRequest struct {
	UserID string `json:"user_id" binding:"required,uuid4"`
}