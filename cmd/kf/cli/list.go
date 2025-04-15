package cli

import (
	"fmt"

	"github.com/joakimen/kf/pkg/kf"
	"github.com/urfave/cli/v2"
)

func newListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List known files",
		Action: func(c *cli.Context) error {
			files, err := kf.List()
			if err != nil {
				fmt.Println("error reading configuration file:", err)
				return err
			}
			for _, file := range files {
				fmt.Println(file)
			}
			return nil
		},
	}
}
