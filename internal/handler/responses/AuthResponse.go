package responses

type AuthResponse struct {
	Token    string `json:"token"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}
