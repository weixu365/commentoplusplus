package util

import (
	"github.com/russross/blackfriday"
)

func MarkdownToHtml(markdown string) string {
	unsafe := blackfriday.Markdown([]byte(markdown), renderer, extensions)
	return string(policy.SanitizeBytes(unsafe))
}
