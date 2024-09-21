package slice

import (
	"slices"
)

func Unique(fileLines []string) []string {
	slices.Sort(fileLines)
	return slices.Compact(fileLines)
}

func Exists(needle string, haystack []string) bool {
	for _, straw := range haystack {
		if needle == straw {
			return true
		}
	}
	return false
}
