package slice

import (
	"strings"
)

func Unique(elements []string) []string {
	seen := make(map[string]struct{}, len(elements))
	result := make([]string, 0, len(elements))
	for _, e := range elements {
		if _, ok := seen[e]; !ok {
			seen[e] = struct{}{}
			result = append(result, e)
		}
	}
	return result
}

func TrimWhitespace(elements []string) []string {
	var trimmedElements []string
	for _, element := range elements {
		trimmedElements = append(trimmedElements, strings.TrimSpace(element))
	}
	return trimmedElements
}

func Exists(needle string, haystack []string) bool {
	for _, straw := range haystack {
		if needle == straw {
			return true
		}
	}
	return false
}
