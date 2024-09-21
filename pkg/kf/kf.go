package kf

import (
	"github.com/joakimen/kf/pkg/slice"
	"github.com/joakimen/kf/pkg/userconfig"
)

func Add(knownFile string) error {
	userConfigLines, err := userconfig.ReadUserConfig()
	if err != nil {
		return err
	}

	if slice.Exists(knownFile, userConfigLines) {
		return userconfig.ErrEntryAlreadyExists
	}

	return userconfig.WriteUserConfig(append(userConfigLines, knownFile))
}

func Forget(knownFile string) (bool, error) {
	userConfigLines, err := userconfig.ReadUserConfig()
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

	err = userconfig.WriteUserConfig(linesToKeep)
	if err != nil {
		return removedMatchingLine, err
	}

	return removedMatchingLine, nil
}

func List() ([]string, error) {
	configFileLines, err := userconfig.ReadUserConfig()
	if err != nil {
		return nil, err
	}
	return configFileLines, nil
}

func Config() (string, error) {
	configFilePath, err := userconfig.GetUserConfigPath()
	if err != nil {
		return "", err
	}
	return configFilePath, nil
}
