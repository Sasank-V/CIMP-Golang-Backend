package utils

import "github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"

func IsValidStatus(status string) bool {
	switch schemas.Status(status) {
	case schemas.Approved, schemas.Pending, schemas.Rejected:
		return true
	default:
		return false
	}
}
