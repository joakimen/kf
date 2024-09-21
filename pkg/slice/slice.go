package slice

import (
	"slices"
	"strings"
)

func Unique(elements []string) []string {
	slices.Sort(elements)
	return slices.Compact(elements)
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
