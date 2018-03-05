package user

type BoundAccountsResponce struct {
	Accounts map[string]string `json:"accounts" binding:"required"`
}

type BoundAccountDeleteRequest struct {
	Resource string `json:"resource" binding:"required"`
}
