package cli

import (
	"fmt"

	"github.com/joakimen/kf/pkg/kf"
	"github.com/urfave/cli/v2"
)

func newAddCmd(getenv func(string) string) *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a file to the list of known files",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return cli.ShowCommandHelp(c, c.Command.Name)
			}

			fileToAdd := c.Args().First()
			if err := kf.Add(fileToAdd, getenv); err != nil {
				return fmt.Errorf("failed when adding to list of known files: %w", err)
			}
			fmt.Println("Added", fileToAdd)
			return nil
		},
	}
}
