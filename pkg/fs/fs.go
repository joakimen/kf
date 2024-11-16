package fs

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type LineReaderWriter interface {
	ReadLines(string) ([]string, error)
	WriteLines([]string) error
}

func ReadLines(filename string) ([]string, error) {
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

func WriteLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(file)

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func SanitizeFilePath(inputPath string) (string, error) {
	var path string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	curDirAbs, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// shorten curdir with ~ if curdir is in the home dir
	curDir := strings.Replace(curDirAbs, homeDir, "~", 1)
	switch {
	case strings.HasPrefix(inputPath, "/"):
		// absolute path
		if strings.HasPrefix(inputPath, homeDir) {
			relPath := inputPath[len(homeDir):]
			path = "~" + relPath
		} else {
			path = inputPath
		}
	case strings.HasPrefix(inputPath, "~"):
		path = inputPath
	default:
		path = filepath.Join(curDir, inputPath)
	}
	return filepath.Clean(path), nil
}
