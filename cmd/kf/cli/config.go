package cli

import (
	"fmt"

	"github.com/joakimen/kf/pkg/kf"
	"github.com/urfave/cli/v2"
)

func newConfigCmd() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Print configuration file path",
		Action: func(c *cli.Context) error {
			config, err := kf.Config()
			if err != nil {
				fmt.Println("error reading configuration file:", err)
				return err
			}
			fmt.Println(config)
			return nil
		},
	}
}
