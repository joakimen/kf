package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joakimen/kf/pkg/fs"
)

func configFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "kf", "config"), nil
}

func AddEntry(line string) error {
	configFilePath, err := configFilePath()
	if err != nil {
		return err
	}
	return fs.AppendToFile(configFilePath, line)
}

func RemoveEntry(line string) ([]fs.Line, error) {
	configFilePath, err := configFilePath()
	if err != nil {
		return nil, err
	}
	return fs.RemoveMatchingLines(configFilePath, line)
}

func ReadConfigFile() ([]string, error) {
	configFilePath, err := configFilePath()
	if err != nil {
		return nil, err
	}
	fileLines, err := fs.ReadFileLines(configFilePath)
	if err != nil {
		return nil, err
	}

	if len(fileLines) == 0 {
		return nil, errors.New("no files found in configuration file")
	}
	return fileLines, nil
}
