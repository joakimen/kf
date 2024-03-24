package fs

import (
	"bufio"
	"cmp"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const defaultEditor = "nvim"

// ReadFileLines reads a file and returns its lines as a slice of slice.
func ReadFileLines(filename string) ([]string, error) {
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

// IsValidFile returns true if the file exists and is a regular file
func IsValidFile(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// GetEditorName returns the name of the editor to use for editing files.
func GetEditorName() string {
	osEditor := os.Getenv("EDITOR")
	return cmp.Or(osEditor, defaultEditor)
}

// EditFile opens the file at the given path in the user's preferred editor.
func EditFile(editor string, path string) error {
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// RealPath expands home and returns the absolute path of a file
func RealPath(path string) (string, error) {
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

func AppendToFile(filename string, line string) error {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// Append the line to the file
	if _, err := file.WriteString(line + "\n"); err != nil {
		log.Fatal(err)
	}
	return nil
}
