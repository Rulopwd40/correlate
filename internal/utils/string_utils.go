package utils

import "strings"

func ReplacePlaceholders(input string, values map[string]string) string {
	for key, value := range values {
		placeholder := "{{" + key + "}}"
		input = strings.ReplaceAll(input, placeholder, value)
	}
	return input
}
