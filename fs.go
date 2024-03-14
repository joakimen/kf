package main

import (
	"bufio"
	"cmp"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	defaultEditor = "nvim"
)

// readFileLines reads a file and returns its lines as a slice of strings.
func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// fileExists returns true if the file exists and false otherwise.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// getEditorName returns the name of the editor to use for editing files.
func getEditorName() string {
	osEditor := os.Getenv("EDITOR")
	return cmp.Or(osEditor, defaultEditor)
}

// editFile opens the file at the given path in the user's preferred editor.
func editFile(editor string, path string) error {
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// expandHome expands a path that starts with ~ to the user's home directory.
func expandHome(path string) (string, error) {

	if path == "" {
		return "", errors.New("input path cannot be blank")
	}

	if path == "~" {
		return os.UserHomeDir()
	}

	if len(path) > 1 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
