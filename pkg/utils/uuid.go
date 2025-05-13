package utils

import (
	"github.com/google/uuid"
	"strings"
)

func StdUUID() string {
	return uuid.New().String()
}
func CompactUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
