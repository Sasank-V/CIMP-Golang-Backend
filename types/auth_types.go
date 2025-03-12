package types

type TokenPayload struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IsLead bool   `json:"is_lead"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
