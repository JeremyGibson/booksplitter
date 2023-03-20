package internal

import (
	"regexp"
	"strings"
)

func normalizeFileName(s string) string {
	newname := strings.ToLower(s)
	newname = regexp.MustCompile(`[^a-z0-9_\-\.]+`).ReplaceAllString(newname, "_")
	newname = regexp.MustCompile(`_{2,}`).ReplaceAllString(newname, "_")
	newname = regexp.MustCompile(`^_|_$`).ReplaceAllString(newname, "")
	return newname
}
