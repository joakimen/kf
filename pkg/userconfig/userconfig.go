package userconfig

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/joakimen/kf/pkg/slice"

	"github.com/joakimen/kf/pkg/fs"
)

func sanitizeUserConfig(lines []string) []string {
	result := slice.Unique(slice.TrimWhitespace(lines))
	sort.Strings(result)
	return result
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
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return fileLines, nil
}
