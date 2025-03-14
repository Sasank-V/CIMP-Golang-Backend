package types

import "github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"

type FullContribution struct {
	Contribution   schemas.Contribution `json:"contribution"`
	ClubName       string               `json:"club_name"`
	DepartmentName string               `json:"department_name"`
}

type AddContributionResponse struct {
	Message        string `json:"message"`
	ContributionID string `json:"cont_id"`
}
