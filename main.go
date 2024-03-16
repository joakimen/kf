package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configFilePath := filepath.Join(homeDir, ".config", "kf", "config")
	fileLines, err := readFileLines(configFilePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "couldn't read configuration file: %v\n", err)
		os.Exit(1)
	}

	uniqueFileLines := unique(fileLines)
	var knownFiles []string
	for _, line := range uniqueFileLines {

		if strings.TrimSpace(line) == "" {
			continue
		}

		absPath, err := expandHome(line)
		if err != nil {
			panic(fmt.Errorf("error expanding home directory: %w", err))

		}

		if !fileExists(absPath) {
			continue
		}

		knownFiles = append(knownFiles, absPath)
	}

	// select one known file
	selectedFile := selectFile(knownFiles)
	fmt.Println(selectedFile)
	err = editFile(getEditorName(), selectedFile)
	if err != nil {
		panic(err)
	}
}

func unique(fileLines []string) []string {
	slices.Sort(fileLines)
	var uniqueFileLines = slices.Compact(fileLines)
	return uniqueFileLines
}
