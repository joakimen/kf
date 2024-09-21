package userconfig

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/joakimen/kf/pkg/slice"

	"github.com/joakimen/kf/pkg/fs"
)

var ErrEntryAlreadyExists = errors.New("entry already exists in configuration file")

func Add(knownFile string) error {
	userConfigLines, err := readUserConfig()
	if err != nil {
		return err
	}

	if slice.Exists(knownFile, userConfigLines) {
		return ErrEntryAlreadyExists
	}

	return writeConfigFile(append(userConfigLines, knownFile))
}

func Forget(knownFile string) (bool, error) {
	userConfigLines, err := readUserConfig()
	if err != nil {
		return false, err
	}

	var linesToKeep []string
	removedMatchingLine := false
	for _, userConfigLine := range userConfigLines {
		if knownFile == userConfigLine {
			removedMatchingLine = true
			continue
		} else {
			linesToKeep = append(linesToKeep, userConfigLine)
		}
	}

	err = writeConfigFile(linesToKeep)
	if err != nil {
		return removedMatchingLine, err
	}

	return removedMatchingLine, nil
}

func List() ([]string, error) {
	configFileLines, err := readUserConfig()
	if err != nil {
		return nil, err
	}
	return configFileLines, nil
}

func sanitizeUserConfig(lines []string) []string {
	uniqueLines := slice.Unique(lines)
	var knownFiles []string
	for _, line := range uniqueLines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		knownFiles = append(knownFiles, trimmedLine)
	}
	return knownFiles
}

func getUserConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "kf", "config"), nil
}

func writeConfigFile(lines []string) error {
	configFilePath, err := getUserConfigFilePath()
	if err != nil {
		return err
	}
	return fs.WriteLines(configFilePath, sanitizeUserConfig(lines))
}

func readUserConfig() ([]string, error) {
	configFilePath, err := getUserConfigFilePath()
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
