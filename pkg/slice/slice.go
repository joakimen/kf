package slice

import (
	"slices"
	"strings"

	"github.com/joakimen/kf/pkg/fs"
)

func unique(fileLines []string) []string {
	slices.Sort(fileLines)
	uniqueFileLines := slices.Compact(fileLines)
	return uniqueFileLines
}

// SanitizeFileSlice removes empty lines, expands home directories in paths, and returns only
// files that exist and are regular files.
func SanitizeFileSlice(lines []string) ([]string, error) {
	uniqueFileLines := unique(lines)
	var knownFiles []string
	for _, line := range uniqueFileLines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		absPath, err := fs.RealPath(line)
		if err != nil {
			return nil, err
		}

		if !fs.IsValidFile(absPath) {
			continue
		}

		knownFiles = append(knownFiles, absPath)
	}
	return knownFiles, nil
}
