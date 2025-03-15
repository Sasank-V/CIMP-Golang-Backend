package types

type UserSignUpInfo struct {
	RegNumber string `json:"reg_number"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserLoginInfo struct {
	RegNumber string `json:"reg_number"`
	Password  string `json:"password"`
}

type TokenPayload struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IsLead bool   `json:"is_lead"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
