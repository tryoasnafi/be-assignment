package auth_supertokens

type AuthResponse struct {
	Status string `json:"status"`
	User   User   `json:"user"`
}

type User struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	TimeJoined int64    `json:"timeJoined"`
	TenantIDS  []string `json:"tenantIds"`
}

type AuthRequest struct {
	FormFields [2]struct {
		ID    string `json:"id" example:"email"`
		Value string `json:"value" example:"email@mail.com"`
	}
}
