package main

import("strings")

func cleanInput(text string) []string {
	var result []string
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return result
	}
	parts := strings.Fields(trimmed)
	for _, part := range parts {
		result = append(result, strings.ToLower(part))
	}
	return result
}