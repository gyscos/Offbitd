package main

import (
	"strings"
)

func sanify(text string) string {
	text = strings.Replace(text, "/", "_", -1)
	text = strings.Replace(text, " ", "_", -1)

	if len(text) > maxFileLength {
		text = text[:maxFileLength]
	}
	return text
}
