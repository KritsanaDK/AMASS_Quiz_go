
package utils

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func GenUUID() string {
	out := uuid.New().String()
	r := strings.NewReplacer("\n", "", "-", "")
	return r.Replace(string(out))
}

func IsValidUUID(uuid string) bool {
	// Check if empty
	if uuid == "" {
		return false
	}

	// Pattern 1: Standard UUID format (8-4-4-4-12)
	// Example: 550e8400-e29b-41d4-a716-446655440000
	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

	return uuidPattern.MatchString(uuid)
}

