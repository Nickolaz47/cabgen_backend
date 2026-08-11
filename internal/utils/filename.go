package utils

import "strings"

func SanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		" ", "_", "/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_", "|", "_",
	)
	return strings.Trim(replacer.Replace(strings.TrimSpace(name)), "_")
}
