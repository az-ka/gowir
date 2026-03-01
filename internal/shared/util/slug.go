package util

import (
	"strings"

	"github.com/gosimple/slug"
)

func GenerateSlug(text string) string {
	formatted := strings.ReplaceAll(text, "&", "dan")
	return slug.Make(formatted)
}
