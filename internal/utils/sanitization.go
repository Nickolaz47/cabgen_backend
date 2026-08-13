package utils

import (
	"strings"

	sanitize "github.com/mrz1836/go-sanitize"
)

func SanitizeQuery(s string) string {
	s = strings.TrimSpace(s)
	return sanitize.SingleLine(s)
}
