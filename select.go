package main

import (
	"fmt"
	fz "github.com/ktr0731/go-fuzzyfinder"
	"os"
)

// Select a known file from a list of known files using fuzzy matching
func selectFile(files []string) string {

	renderFunc := func(selectedIndex int) string {
		return files[selectedIndex]
	}

	previewFunc := func(selectedIndex, width, height int) string {
		if selectedIndex == -1 {
			return ""
		}
		file := files[selectedIndex]
		contents, err := os.ReadFile(file)
		if err != nil {
			return fmt.Sprintf("Error reading file: %s", err)
		}

		colorizedContents, err := Colorize(file, string(contents))
		if err != nil {
			return fmt.Sprintf("Error colorizing file: %s", err)
		}
		return colorizedContents
	}

	idx, err := fz.Find(files, renderFunc, fz.WithPreviewWindow(previewFunc))
	if err != nil {
		fmt.Println("no file selected")
		os.Exit(0)
	}

	return files[idx]
}
