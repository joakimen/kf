package cli

import (
	"fmt"

	"github.com/joakimen/kf/pkg/kf"
	"github.com/urfave/cli/v2"
)

func newForgetCmd() *cli.Command {
	return &cli.Command{
		Name:  "forget",
		Usage: "Remove a file from the list of known files",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return cli.ShowCommandHelp(c, c.Command.Name)
			}

			fileToRemove := c.Args().First()
			removed, err := kf.Forget(fileToRemove)
			if err != nil {
				return fmt.Errorf("failed when removing from list of known files: %w", err)
			}
			if removed {
				fmt.Println("Removed", fileToRemove)
			} else {
				fmt.Println("No matching entry found")
			}
			return nil
		},
	}
}
