package validation

import "strings"

// IsValidEmail checks if email matches a rudimentary format check
func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
