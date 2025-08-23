package utils

import (
	"strings"

	"github.com/google/uuid"
)

func StdUUID() string {
	return uuid.New().String()
}
func CompactUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
