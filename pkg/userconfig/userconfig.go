package userconfig

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joakimen/kf/pkg/slice"

	"github.com/joakimen/kf/pkg/fs"
)

func sanitizeUserConfig(lines []string) []string {
	uniqueLines := slice.Unique(lines)
	trimmedElements := slice.TrimWhitespace(uniqueLines)
	return trimmedElements
}

func GetUserConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "kf", "config"), nil
}

func WriteUserConfig(lines []string) error {
	configFilePath, err := GetUserConfigPath()
	if err != nil {
		return err
	}
	return fs.WriteLines(configFilePath, sanitizeUserConfig(lines))
}

func ReadUserConfig() ([]string, error) {
	configFilePath, err := GetUserConfigPath()
	if err != nil {
		return nil, err
	}
	fileLines, err := fs.ReadLines(configFilePath)
	if err != nil {
		return nil, err
	}

	if len(fileLines) == 0 {
		return nil, errors.New("no files found in configuration file")
	}
	return fileLines, nil
}
