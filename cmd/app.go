package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joakimen/kf/pkg/kf"
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
			{
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
			},
			{
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
			},
			{
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
			},
			{
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
			},
		},
	}
}
