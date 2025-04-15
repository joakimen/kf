package cli

import (
	"github.com/urfave/cli/v2"
)

func NewApp(getenv func(string) string) cli.App {
	return cli.App{
		Name:  "kf",
		Usage: "Manages known files",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose output",
			},
		},
		Commands: []*cli.Command{
			newConfigCmd(),
			newAddCmd(getenv),
			newListCmd(),
			newEditCmd(),
			newForgetCmd(),
		},
	}
}
