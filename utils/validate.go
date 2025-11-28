package utils

import "strings"

func ValidEmail(email string) bool {
	return strings.Contains(email, "@") && len(email) > 3
}
