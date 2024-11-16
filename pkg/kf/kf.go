package kf

import (
	"errors"

	"github.com/joakimen/kf/pkg/fs"
	"github.com/joakimen/kf/pkg/slice"
	"github.com/joakimen/kf/pkg/userconfig"
)

var (
	ErrEntryAlreadyExists    = errors.New("entry already exists in configuration file")
	ErrCannotReadUserConfig  = errors.New("error reading configuration file")
	ErrCannotWriteUserConfig = errors.New("error writing configuration file")
)

func Add(inputKnownFile string, getenv func(string) string) error {
	userConfigLines, err := userconfig.ReadUserConfig()
	if err != nil {
		return errors.Join(ErrCannotReadUserConfig, err)
	}

	knownFile, err := fs.SanitizeFilePath(inputKnownFile, getenv)
	if err != nil {
		return err
	}

	if slice.Exists(knownFile, userConfigLines) {
		return ErrEntryAlreadyExists
	}
	return userconfig.WriteUserConfig(append(userConfigLines, knownFile))
}

func Forget(knownFile string) (bool, error) {
	userConfigLines, err := userconfig.ReadUserConfig()
	if err != nil {
		return false, errors.Join(ErrCannotReadUserConfig, err)
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

	err = userconfig.WriteUserConfig(linesToKeep)
	if err != nil {
		return removedMatchingLine, errors.Join(ErrCannotWriteUserConfig, err)
	}

	return removedMatchingLine, nil
}

func List() ([]string, error) {
	configFileLines, err := userconfig.ReadUserConfig()
	if err != nil {
		return nil, errors.Join(ErrCannotReadUserConfig, err)
	}
	return configFileLines, nil
}

func Config() (string, error) {
	configFilePath, err := userconfig.GetUserConfigPath()
	if err != nil {
		return "", errors.Join(ErrCannotReadUserConfig, err)
	}
	return configFilePath, nil
}
