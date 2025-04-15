package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/joakimen/kf/pkg/userconfig"
	"github.com/urfave/cli/v2"
)

func newEditCmd() *cli.Command {
	return &cli.Command{
		Name:        "edit",
		Description: "Edit configuration file",
		Action: func(c *cli.Context) error {
			osEditor := os.Getenv("EDITOR")
			if osEditor == "" {
				return errors.New("$EDITOR env var is not set")
			}

			configFilePath, err := userconfig.GetUserConfigPath()
			if err != nil {
				return fmt.Errorf("error getting configuration file path: %w", err)
			}

			editCmd := exec.Command(osEditor, configFilePath)
			editCmd.Stdout = os.Stdout
			editCmd.Stderr = os.Stderr
			return editCmd.Run()
		},
	}
}
