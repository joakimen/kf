package fs

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Line struct {
	Number int
	Text   string
}

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

func RemoveMatchingLines(filename string, pattern string) ([]Line, error) {
	lines, err := ReadFileLines(filename)
	if err != nil {
		return nil, err
	}

	var linesToKeep []string
	var removedLines []Line
	for lineNum, line := range lines {
		realPath, err := RealPath(line)
		if err != nil {
			return nil, err
		}
		if pattern == realPath {
			removedLines = append(removedLines, Line{Number: lineNum, Text: line})
			continue
		}
		linesToKeep = append(linesToKeep, line)
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range linesToKeep {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return nil, err
		}
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	return removedLines, nil
}

// IsValidFile returns true if the file exists and is a regular file
func IsValidFile(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
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
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
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
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
