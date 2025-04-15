package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/joakimen/kf/pkg/kf"
	"github.com/urfave/cli/v2"
)

func newAddCmd(getenv func(string) string) *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a file to the list of known files",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				log.Fatal("No CLI arguments detected!")
			}

			fileToAdd := c.Args().First()
			err := kf.Add(fileToAdd, getenv)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed when adding to list of known files: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Added", fileToAdd)
			return nil
		},
	}
}
