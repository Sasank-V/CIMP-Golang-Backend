package types

import "github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"

type GetUserResponse struct {
	Message string       `json:"message"`
	User    schemas.User `json:"user,omitempty"`
}

type GetUserRequestsResponse struct {
	Message       string             `json:"message"`
	Contributions []FullContribution `json:"requests"`
}
