package types

import "github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"

type GetUserResponse struct {
	Message string       `json:"message"`
	User    schemas.User `json:"user,omitempty"`
}

type GetUserContributionsResponse struct {
	Message       string             `json:"message"`
	Contributions []FullContribution `json:"contributions"`
}

type GetLeadUserRequestsResponse struct {
	Message  string             `json:"message"`
	Requests []FullContribution `json:"requests"`
}
