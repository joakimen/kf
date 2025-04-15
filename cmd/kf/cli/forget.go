package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/joakimen/kf/pkg/kf"
	"github.com/urfave/cli/v2"
)

func newForgetCmd() *cli.Command {
	return &cli.Command{
		Name:  "forget",
		Usage: "Remove a file from the list of known files",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				log.Fatal("No CLI arguments detected!")
			}

			fileToRemove := c.Args().First()
			removed, err := kf.Forget(fileToRemove)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed when removing from list of known files: %v\n", err)
				os.Exit(1)
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
