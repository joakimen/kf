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
	if err := os.MkdirAll(filepath.Dir(filePath), 0o700); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(filePath), ".kf-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		if tmp != nil {
			tmp.Close()
			os.Remove(tmpPath)
		}
	}()

	w := bufio.NewWriter(tmp)
	for _, line := range lines {
		if _, err := w.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	if err := w.Flush(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	tmp = nil
	return os.Rename(tmpPath, filePath)
}

func SanitizeFilePath(inputPath string, getenv func(string) string) (string, error) {
	var path string

	homeDir := getenv("HOME")
	if homeDir == "" {
		return "", errors.New("HOME not set")
	}

	curDirAbs := getenv("PWD")
	if curDirAbs == "" {
		return "", errors.New("PWD not set")
	}

	// shorten curdir with ~ if curdir is in the home dir
	curDir := curDirAbs
	if strings.HasPrefix(curDirAbs, homeDir) {
		curDir = "~" + curDirAbs[len(homeDir):]
	}
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
